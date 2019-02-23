package runner

import (
	"context"
	"fmt"

	"github.com/alecthomas/kingpin"
	kitlog "github.com/go-kit/kit/log"
	workloadsv1alpha1 "github.com/gocardless/theatre/pkg/apis/workloads/v1alpha1"
	"github.com/gocardless/theatre/pkg/client/clientset/versioned"
	theatre "github.com/gocardless/theatre/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	CLI = kingpin.New("consoles", "Manages theatre consoles")

	create         = CLI.Command("create", "Creates a new console given a template")
	createSelector = create.Flag("selector", "Selector that matches console template").Required().String()
	createTimeout  = create.Flag("timeout", "Timeout for the new console").Duration()
	createReason   = create.Flag("reason", "Reason for creating console").String()
	createCommand  = create.Arg("command", "Command to run in console").Strings()

	list   = CLI.Command("list", "Lists cluster consoles")
	attach = CLI.Command("attach", "Attaches to existing console")
)

func CLIRun(ctx context.Context, logger kitlog.Logger, config *rest.Config, args []string) error {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	theatreClient, err := theatre.NewForConfig(config)
	if err != nil {
		return err
	}

	runner := New(client, theatreClient)

	switch kingpin.MustParse(CLI.Parse(args)) {
	case create.FullCommand():
		tpl, err := runner.FindTemplateBySelector(metav1.NamespaceAll, *createSelector)
		if err != nil {
			return err
		}

		opt := Options{Cmd: *createCommand, Timeout: int(createTimeout.Seconds()), Reason: *createReason}
		csl, err := runner.Create(tpl.Namespace, *tpl, opt)
		if err != nil {
			return nil
		}

		csl, err = runner.WaitUntilReady(ctx, *csl)
		if err != nil {
			return nil
		}

		pod, err := runner.GetAttachablePod(csl)
		if err != nil {
			return nil
		}

		logger.Log("pod", pod.Name, "msg", "console pod created")
	case list.FullComment():

	}

	return nil
}

// Runner is responsible for managing the lifecycle of a console
type Runner struct {
	kubeClient    kubernetes.Interface
	theatreClient versioned.Interface
}

// Options defines the parameters that can be set upon a new console
type Options struct {
	Cmd     []string
	Timeout int
	Reason  string

	// TODO: For now we assume that all consoles are interactive, i.e. we setup a TTY on
	// them when spawning them. This does not enforce a requirement to attach to the console
	// though.
	// Later on we may need to implement non-interactive consoles, for processes which
	// expect a TTY to not be present?
	// However with these types of consoles it will not be possible to send input to them
	// when reattaching, e.g. attempting to send a SIGINT to cancel the running process.
	// Interactive bool
}

// New builds a runner
func New(client kubernetes.Interface, theatreClient versioned.Interface) *Runner {
	return &Runner{
		kubeClient:    client,
		theatreClient: theatreClient,
	}
}

// Create builds a console according to the supplied options and submits it to the API
func (c *Runner) Create(namespace string, template workloadsv1alpha1.ConsoleTemplate, opts Options) (*workloadsv1alpha1.Console, error) {
	csl := &workloadsv1alpha1.Console{
		ObjectMeta: metav1.ObjectMeta{
			// Let Kubernetes generate a unique name
			GenerateName: template.Name + "-",
			Labels:       labels.Merge(labels.Set{}, template.Labels),
		},
		Spec: workloadsv1alpha1.ConsoleSpec{
			ConsoleTemplateRef: corev1.LocalObjectReference{Name: template.Name},
			// If the flag is not provided then the value will default to 0. The controller
			// should detect this and apply the default timeout that is defined in the template.
			TimeoutSeconds: opts.Timeout,
			Command:        opts.Cmd,
			Reason:         opts.Reason,
		},
	}

	return c.theatreClient.WorkloadsV1alpha1().Consoles(namespace).Create(csl)
}

// FindTemplateBySelector will search for a template matching the given label
// selector and return errors if none or multiple are found (when the selector
// is too broad)
func (c *Runner) FindTemplateBySelector(namespace string, labelSelector string) (*workloadsv1alpha1.ConsoleTemplate, error) {
	client := c.theatreClient.WorkloadsV1alpha1().ConsoleTemplates(namespace)

	templates, err := client.List(
		metav1.ListOptions{
			LabelSelector: labelSelector,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list consoles templates")
	}

	if len(templates.Items) != 1 {
		identifiers := []string{}
		for _, item := range templates.Items {
			identifiers = append(identifiers, item.Namespace+"/"+item.Name)
		}

		return nil, errors.Errorf(
			"expected to discover 1 console template, but actually found: %s",
			identifiers,
		)
	}

	template := templates.Items[0]

	return &template, nil
}

func (c *Runner) FindConsoleByName(namespace, name string) (*workloadsv1alpha1.Console, error) {
	// We must List then filter the slice instead of calling Get(name), otherwise
	// the real Kubernetes client will return the following error when namespace
	// is empty: "an empty namespace may not be set when a resource name is
	// provided".
	// The fake clientset generated by client-gen will not replicate this error in
	// unit tests.
	allConsolesInNamespace, err := c.theatreClient.WorkloadsV1alpha1().Consoles(namespace).
		List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var matchingConsoles []workloadsv1alpha1.Console
	for _, console := range allConsolesInNamespace.Items {
		if console.Name == name {
			matchingConsoles = append(matchingConsoles, console)
		}
	}

	if len(matchingConsoles) == 0 {
		return nil, fmt.Errorf("no consoles found with name: %s", name)
	}
	if len(matchingConsoles) > 1 {
		return nil, fmt.Errorf("too many consoles found with name: %s, please specify namespace", name)
	}

	return &matchingConsoles[0], nil
}

func (c *Runner) ListConsolesByLabelsAndUser(namespace, username, labelSelector string) ([]workloadsv1alpha1.Console, error) {
	// We cannot use a FieldSelector on spec.user in conjunction with the
	// LabelSelector for CRD types like Console. The error message "field label
	// not supported: spec.user" is returned by the real Kubernetes client.
	// See https://github.com/kubernetes/kubernetes/issues/53459.
	csls, err := c.theatreClient.WorkloadsV1alpha1().Consoles(namespace).
		List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}

	var filtered []workloadsv1alpha1.Console
	for _, csl := range csls.Items {
		if username == "" || csl.Spec.User == username {
			filtered = append(filtered, csl)
		}
	}
	return filtered, err
}

// WaitUntilReady will block until the console reaches a phase that indicates
// that it's ready to be attached to, or has failed.
// It will then block until an associated RoleBinding exists that contains the
// console user in its subject list. This RoleBinding gives the console user
// permission to attach to the pod.
func (c *Runner) WaitUntilReady(ctx context.Context, createdCsl workloadsv1alpha1.Console) (*workloadsv1alpha1.Console, error) {
	csl, err := c.waitForConsole(ctx, createdCsl)
	if err != nil {
		return nil, err
	}

	if err := c.waitForRoleBinding(ctx, csl); err != nil {
		return nil, err
	}

	return csl, nil
}

func (c *Runner) waitForConsole(ctx context.Context, createdCsl workloadsv1alpha1.Console) (*workloadsv1alpha1.Console, error) {
	isRunning := func(csl *workloadsv1alpha1.Console) bool {
		return csl != nil && csl.Status.Phase == workloadsv1alpha1.ConsoleRunning
	}
	isStopped := func(csl *workloadsv1alpha1.Console) bool {
		return csl != nil && csl.Status.Phase == workloadsv1alpha1.ConsoleStopped
	}

	listOptions := metav1.SingleObject(createdCsl.ObjectMeta)
	client := c.theatreClient.WorkloadsV1alpha1().Consoles(createdCsl.Namespace)

	w, err := client.Watch(listOptions)
	if err != nil {
		return nil, errors.Wrap(err, "error watching console")
	}

	// Get the console, because watch will only give us an event when something
	// is changed, and the phase could have already stabilised before the watch
	// is set up.
	csl, err := client.Get(createdCsl.Name, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, errors.Wrap(err, "error retrieving console")
	}

	// If the console is already running then there's nothing to do
	if isRunning(csl) {
		return csl, nil
	}
	if isStopped(csl) {
		return nil, fmt.Errorf("console is Stopped")
	}

	status := w.ResultChan()
	defer w.Stop()

	for {
		select {
		case event, ok := <-status:
			// If our channel is closed, exit with error, as we'll otherwise assume
			// we were successful when we never reached this state.
			if !ok {
				return nil, fmt.Errorf("watch channel closed")
			}

			csl = event.Object.(*workloadsv1alpha1.Console)
			if isRunning(csl) {
				return csl, nil
			}
			if isStopped(csl) {
				return nil, fmt.Errorf("console is Stopped")
			}
		case <-ctx.Done():
			if csl == nil {
				return nil, errors.Wrap(ctx.Err(), "console not found")
			}
			return nil, errors.Wrap(ctx.Err(), fmt.Sprintf(
				"console's last phase was: '%v'", csl.Status.Phase),
			)
		}
	}
}

func (c *Runner) waitForRoleBinding(ctx context.Context, csl *workloadsv1alpha1.Console) error {
	rbClient := c.kubeClient.RbacV1().RoleBindings(csl.Namespace)
	watcher, err := rbClient.Watch(metav1.ListOptions{FieldSelector: "metadata.name=" + csl.Name})
	if err != nil {
		return errors.Wrap(err, "error watching rolebindings")
	}
	defer watcher.Stop()

	// The Console controller might have already created a DirectoryRoleBinding
	// and the DirectoryRoleBinding controller might have created the RoleBinding
	// and updated its subject list by this point. If so, we are already done, and
	// might never receive another event from our RoleBinding Watcher, causing the
	// subsequent loop would block forever.
	// If the associated RoleBinding exists and has the console user in its
	// subject list, return early.
	rb, err := rbClient.Get(csl.Name, metav1.GetOptions{})
	if err == nil && rbHasSubject(rb, csl.Spec.User) {
		return nil
	}

	rbEvents := watcher.ResultChan()
	for {
		select {
		case rbEvent, ok := <-rbEvents:
			if !ok {
				return errors.New("rolebinding event watcher channel closed")
			}

			rb := rbEvent.Object.(*rbacv1.RoleBinding)
			if rbHasSubject(rb, csl.Spec.User) {
				return nil
			}

			continue

		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "waiting for rolebinding interrupted")
		}
	}
}

func rbHasSubject(rb *rbacv1.RoleBinding, subjectName string) bool {
	for _, subject := range rb.Subjects {
		if subject.Name == subjectName {
			return true
		}
	}
	return false
}

// GetAttachablePod returns an attachable pod for the given console
func (c *Runner) GetAttachablePod(csl *workloadsv1alpha1.Console) (*corev1.Pod, error) {
	pod, err := c.kubeClient.CoreV1().Pods(csl.Namespace).Get(csl.Status.PodName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	for _, c := range pod.Spec.Containers {
		if c.TTY {
			return pod, nil
		}
	}

	return nil, errors.New("no attachable pod found")
}
