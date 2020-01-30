package controllers

func (r *GuestBookReconciler) desiredService(book webappv1.GuestBook) (corev1.Service, error) {
    svc := corev1.Service{
        TypeMeta: metav1.TypeMeta{
            APIVersion: corev1.SchemeGroupVersion.String(), Kind: "Service",
        },
        ObjectMeta: metav1.ObjectMeta{Name: book.Name, Namespace: book.Namespace},
        Spec: corev1.ServiceSpec{
            Ports: []corev1.ServicePort{
                {Name: "http", Port: 8080,
                 Protocol: "TCP", TargetPort: intstr.FromString("http")},
            },
            Selector: map[string]string{"guestbook": book.Name},
            Type:     corev1.ServiceTypeLoadBalancer,
        },
    }

    // always set the controller reference so that we know which object owns this.
    if err := ctrl.SetControllerReference(&book, &svc, r.Scheme); err != nil {
        return svc, err
    }
    return svc, nil
}

func urlForService(svc corev1.Service, port int32) string {
    // unset this if it's not present -- we always want the
    // state to reflect what we observe.
    if len(svc.Status.LoadBalancer.Ingress) == 0 {
        return ""
    }

    host := svc.Status.LoadBalancer.Ingress[0].Hostname
    if host == "" {
        host = svc.Status.LoadBalancer.Ingress[0].IP
    }
    return fmt.Sprintf("http://%s", net.JoinHostPort(host, fmt.Sprintf("%v", port)))
}
