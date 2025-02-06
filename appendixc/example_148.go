// Example 148
// internal/config/manager.go
type ConfigManager struct {
    vaultClient *vault.Client
    k8sClient   *kubernetes.Clientset
    watchChan   chan struct{}
}

func (cm *ConfigManager) WatchConfigMap(ctx context.Context, name string) error {
    watcher, err := cm.k8sClient.CoreV1().ConfigMaps("default").Watch(ctx, metav1.ListOptions{
        FieldSelector: fmt.Sprintf("metadata.name=%s", name),
    })
    if err != nil {
        return fmt.Errorf("watching configmap: %w", err)
    }
    
    go func() {
        for event := range watcher.ResultChan() {
            if event.Type == watch.Modified {
                cm.watchChan <- struct{}{}
            }
        }
    }()
    
    return nil
}

func (cm *ConfigManager) LoadSecrets(ctx context.Context, path string) (map[string]string, error) {
    secret, err := cm.vaultClient.KVv2("secret").Get(ctx, path)
    if err != nil {
        return nil, fmt.Errorf("getting secrets: %w", err)
    }
    
    return secret.Data, nil
}