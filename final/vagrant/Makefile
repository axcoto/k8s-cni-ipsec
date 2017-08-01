.DEFAULT_GOAL := up
KUBECTL := sudo kubectl --kubeconfig /etc/kubernetes/admin.conf
NAMESPACE := --namespace=demo

# shortcut to bring everything up and run depoyment
up: vagrant deploy get_pods

vagrant:
	vagrant up

# Deploy and waiting until all pod are ready
deploy:
	vagrant ssh master -c "$(KUBECTL) create -f /vagrant/echo-demo-pod.yaml"
	@echo "Sleep 10 seconds to wait for pod"

# Attach shell to a pod
shell:
	$(eval Pod := $(shell vagrant ssh master -c "sudo docker ps | grep k8s_echo-demo_" | awk -F'k8s_echo-demo_' '{print $$2}' | awk -F'_' '{print $$1}' | head -n1))
	vagrant ssh master -c "sudo $(KUBECTL) $(NAMESPACE) exec -it $(Pod) -- sh"


# Dev tasks
get_pods:
	vagrant ssh master -c "$(KUBECTL) $(NAMESPACE) get pods"

teardown:
	vagrant ssh master -c "$(KUBECTL) delete -f /vagrant/echo-demo-pod.yaml"

find_ip:
	vagrant ssh master -c "sudo ipsec statusall" | grep 'server.*10.173.0.0.*/32' | awk -F'===' '{print $$2}' | head -n9

reset_kubeadm:
	vagrant ssh master -c 'sudo kubeadm reset'
