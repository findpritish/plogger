package main

import "fmt"

type Target struct {
	Namespace string
	Pod       string
	Container string
}

// GetID returns the ID of the object
func (t *Target) GetID() string {
	return fmt.Sprintf("%s-%s-%s", t.Namespace, t.Pod, t.Container)
}

func Watch(i v1.PodInterface, podFilter *regexp.Regexp, containerFilter *regexp.Regexp, containerExcludeFilter *regexp.Regexp, containerState ContainerState) (chan *Target, chan *Target, error) {
	watcher, err := i.Watch(metav1.ListOptions{Watch: true})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to set up watch")
	}

	added := make(chan *Target)
	removed := make(chan *Target)

	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				if e.Object == nil {
					// Closed because of error
					return
				}

				pod := e.Object.(*corev1.Pod)

				if !podFilter.MatchString(pod.Name) {
					continue
				}

				switch e.Type {
				case watch.Added, watch.Modified:
					var statuses []corev1.ContainerStatus
					statuses = append(statuses, pod.Status.InitContainerStatuses...)
					statuses = append(statuses, pod.Status.ContainerStatuses...)

					for _, c := range statuses {
						if !containerFilter.MatchString(c.Name) {
							continue
						}
						if containerExcludeFilter != nil && containerExcludeFilter.MatchString(c.Name) {
							continue
						}

						if containerState.Match(c.State) {
							added <- &Target{
								Namespace: pod.Namespace,
								Pod:       pod.Name,
								Container: c.Name,
							}
						}
					}
				case watch.Deleted:
					var containers []corev1.Container
					containers = append(containers, pod.Spec.Containers...)
					containers = append(containers, pod.Spec.InitContainers...)

					for _, c := range containers {
						if !containerFilter.MatchString(c.Name) {
							continue
						}
						if containerExcludeFilter != nil && containerExcludeFilter.MatchString(c.Name) {
							continue
						}

						removed <- &Target{
							Namespace: pod.Namespace,
							Pod:       pod.Name,
							Container: c.Name,
						}
					}
				}
			case <-ctx.Done():
				watcher.Stop()
				close(added)
				close(removed)
				return
			}
		}
	}()

	return added, removed, nil
}