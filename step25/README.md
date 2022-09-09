#Example #1
```
def simple_persistentvolumeclaim():
    """Return the Kubernetes config matching the simple-persistentvolumeclaim.yaml manifest."""
    return client.V1PersistentVolumeClaim(
        api_version='v1',
        kind='PersistentVolumeClaim',
        metadata=client.V1ObjectMeta(
            name='my-pvc'
        ),
        spec=client.V1PersistentVolumeClaimSpec(
            access_modes=[
                'ReadWriteMany'
            ],
            resources=client.V1ResourceRequirements(
                requests={
                    'storage': '16Mi'
                }
            )
        )
    ) 
```

#Example #2
```
def get_volume_claim_templates(self) -> Sequence[V1PersistentVolumeClaim]:
        return [
            V1PersistentVolumeClaim(
                metadata=V1ObjectMeta(name=self.get_persistent_volume_name(volume)),
                spec=V1PersistentVolumeClaimSpec(
                    # must be ReadWriteOnce for EBS
                    access_modes=["ReadWriteOnce"],
                    storage_class_name=self.get_storage_class_name(volume),
                    resources=V1ResourceRequirements(
                        requests={"storage": f"{volume['size']}Gi"}
                    ),
                ),
            )
            for volume in self.get_persistent_volumes()
        ] 
```

#Example #3
```
def _createKubePVC(parser_args,CoreV1Api):
    if parser_args.component=="search":
        # {tenantName}-{environment}-search-master-gluster-volume
        if parser_args.envtype == "auth":
            _pvc_name=parser_args.tenant+parser_args.env+"-search-master-"+"volume"
        elif parser_args.envtype == "live":
            _pvc_name=parser_args.tenant+parser_args.env+"-search-repeater-"+"volume"

        # Below config boby still need refine to accept more flexible extension
        metadata = {'name': _pvc_name, 'namespace': parser_args.namespace}
        requests = {'storage': parser_args.storage_size}
        _V1ResourceRequirements=client.V1ResourceRequirements(requests=requests)
        _V1PersistentVolumeClaimSpec=client.V1PersistentVolumeClaimSpec(resources=_V1ResourceRequirements,storage_class_name=parser_args.storage_class,access_modes=['ReadWriteMany'])
        body = client.V1PersistentVolumeClaim(api_version='v1',kind='PersistentVolumeClaim', metadata=metadata,spec=_V1PersistentVolumeClaimSpec,)  # V1PersistentVolumeClaim |

        try:
             api_response = CoreV1Api.create_namespaced_persistent_volume_claim(parser_args.namespace, body)
             print(api_response)
        except ApiException as e:
            print("Exception when calling CoreV1Api->create_namespaced_persistent_volume_claim: %s\n" % e)
    else:
        print("Compoennt %s not need to create persistent volume chain on Kubernetes" %(parser_args.component)) 
```

#Example #4
```
def test_get_volume_claim_templates(self):
        with mock.patch(
            "paasta_tools.kubernetes_tools.KubernetesDeploymentConfig.get_persistent_volumes",
            autospec=True,
        ) as mock_get_persistent_volumes, mock.patch(
            "paasta_tools.kubernetes_tools.KubernetesDeploymentConfig.get_persistent_volume_name",
            autospec=True,
        ) as mock_get_persistent_volume_name, mock.patch(
            "paasta_tools.kubernetes_tools.KubernetesDeploymentConfig.get_storage_class_name",
            autospec=True,
        ) as mock_get_storage_class_name:
            mock_get_persistent_volumes.return_value = [{"size": 20}, {"size": 10}]
            expected = [
                V1PersistentVolumeClaim(
                    metadata=V1ObjectMeta(
                        name=mock_get_persistent_volume_name.return_value
                    ),
                    spec=V1PersistentVolumeClaimSpec(
                        access_modes=["ReadWriteOnce"],
                        storage_class_name=mock_get_storage_class_name.return_value,
                        resources=V1ResourceRequirements(requests={"storage": "10Gi"}),
                    ),
                ),
                V1PersistentVolumeClaim(
                    metadata=V1ObjectMeta(
                        name=mock_get_persistent_volume_name.return_value
                    ),
                    spec=V1PersistentVolumeClaimSpec(
                        access_modes=["ReadWriteOnce"],
                        storage_class_name=mock_get_storage_class_name.return_value,
                        resources=V1ResourceRequirements(requests={"storage": "20Gi"}),
                    ),
                ),
            ]
            ret = self.deployment.get_volume_claim_templates()
            assert expected[0] in ret
            assert expected[1] in ret
            assert len(ret) == 2 
```

#Example #5

```
def _AssumeInsideKfp(
      self,
      namespace='my-namespace',
      pod_name='my-pod-name',
      pod_uid='my-pod-uid',
      pod_service_account_name='my-service-account-name',
      with_pvc=False):
    pod = k8s_client.V1Pod(
        api_version='v1',
        kind='Pod',
        metadata=k8s_client.V1ObjectMeta(
            name=pod_name,
            uid=pod_uid,
        ),
        spec=k8s_client.V1PodSpec(
            containers=[
                k8s_client.V1Container(
                    name='main',
                    volume_mounts=[]),
            ],
            volumes=[]))

    if with_pvc:
      pod.spec.volumes.append(
          k8s_client.V1Volume(
              name='my-volume',
              persistent_volume_claim=k8s_client
              .V1PersistentVolumeClaimVolumeSource(
                  claim_name='my-pvc')))
      pod.spec.containers[0].volume_mounts.append(
          k8s_client.V1VolumeMount(
              name='my-volume',
              mount_path=self._base_dir))

    mock.patch.object(kube_utils, 'is_inside_kfp', return_value=True).start()
    pod.spec.service_account_name = pod_service_account_name
    mock.patch.object(kube_utils, 'get_current_kfp_pod',
                      return_value=pod).start()
    mock.patch.object(kube_utils, 'get_kfp_namespace',
                      return_value=namespace).start()
    if with_pvc:
      (self._mock_core_v1_api.read_namespaced_persistent_volume_claim
       .return_value) = k8s_client.V1PersistentVolumeClaim(
           metadata=k8s_client.V1ObjectMeta(
               name='my-pvc'),
           spec=k8s_client.V1PersistentVolumeClaimSpec(
               access_modes=['ReadWriteMany'])) 
```



#Example #6


```
def create_k8s_nfs_resources(self) -> bool:
        """
        Create NFS resources such as PV and PVC in Kubernetes.
        """
        from kubernetes import client as k8sclient

        pv_name = "nfs-ckpt-pv-{}".format(uuid.uuid4())
        persistent_volume = k8sclient.V1PersistentVolume(
            api_version="v1",
            kind="PersistentVolume",
            metadata=k8sclient.V1ObjectMeta(
                name=pv_name,
                labels={'app': pv_name}
            ),
            spec=k8sclient.V1PersistentVolumeSpec(
                access_modes=["ReadWriteMany"],
                nfs=k8sclient.V1NFSVolumeSource(
                    path=self.params.path,
                    server=self.params.server
                ),
                capacity={'storage': '10Gi'},
                storage_class_name=""
            )
        )
        k8s_api_client = k8sclient.CoreV1Api()
        try:
            k8s_api_client.create_persistent_volume(persistent_volume)
            self.params.pv_name = pv_name
        except k8sclient.rest.ApiException as e:
            print("Got exception: %s\n while creating the NFS PV", e)
            return False

        pvc_name = "nfs-ckpt-pvc-{}".format(uuid.uuid4())
        persistent_volume_claim = k8sclient.V1PersistentVolumeClaim(
            api_version="v1",
            kind="PersistentVolumeClaim",
            metadata=k8sclient.V1ObjectMeta(
                name=pvc_name
            ),
            spec=k8sclient.V1PersistentVolumeClaimSpec(
                access_modes=["ReadWriteMany"],
                resources=k8sclient.V1ResourceRequirements(
                    requests={'storage': '10Gi'}
                ),
                selector=k8sclient.V1LabelSelector(
                    match_labels={'app': self.params.pv_name}
                ),
                storage_class_name=""
            )
        )

        try:
            k8s_api_client.create_namespaced_persistent_volume_claim(self.params.namespace, persistent_volume_claim)
            self.params.pvc_name = pvc_name
        except k8sclient.rest.ApiException as e:
            print("Got exception: %s\n while creating the NFS PVC", e)
            return False

        return True 
```

###
```
from kubernetes import config,client
import os,hashlib,random
import json
#基于token认证
token='xxx'
configuration = client.Configuration()
configuration.host = "https://192.168.1.1:6443" # APISERVER地址
configuration.ssl_ca_cert = "/Users/wuchenyang/code/devops/kubeconfig/ca.crt"  # CA证书
configuration.verify_ssl = True  # 启用证书验证
configuration.api_key = {"authorization": "Bearer " + token}  # 指定Token字符串
client.Configuration.set_default(configuration)
sc_v1=client.StorageV1Api()
apps_v1=client.AppsV1Api()
core_api = client.CoreV1Api()
# for i in apps_v1.list_namespaced_deployment(namespace='default').items:
#     print(i.metadata.name)

print("===========")
for sc in sc_v1.list_storage_class().items:
    print(sc.metadata.name)
    print(sc.provisioner)
    print(sc.reclaim_policy)
    print(sc.volume_binding_mode)
    print(sc.allow_volume_expansion)
    print(sc.metadata.creation_timestamp)
    # print("==========1===========")
    #
    # print("===========2===========")
    # print(sc.metadata.managed_fields)
print("===========")
# for pv in core_api.list_persistent_volume().items:
#     print(pv.metadata.name)
#     print(pv.spec.capacity['storage'])
name='test'
capacity='50Gi'
body=client.V1PersistentVolume(
    api_version="v1",
    kind="PersistentVolume",
    metadata=client.V1ObjectMeta(name=name),
    spec=client.V1PersistentVolumeSpec(
        capacity={'storage': '50Gi'},
        access_modes=["ReadWriteOnce"],
        volume_mode="Filesystem",
        storage_class_name="managed-nfs-storage",
        nfs=client.V1NFSVolumeSource(
            server='192.168.1.1',
            path='/ifs/kubernetes'
        ),

    )
)
core_api.create_persistent_volume(body=body)

name='my-pvc-test'
bodypvc=client.V1PersistentVolumeClaim(
    api_version="v1",
    kind="PersistentVolumeClaim",
    metadata=client.V1ObjectMeta(name=name),
    spec=client.V1PersistentVolumeClaimSpec(
        access_modes=["ReadWriteOnce"],
        storage_class_name='managed-nfs-storage',
        resources=client.V1ResourceRequirements(
            requests={"storage":"40Gi"}
        )
    )
)
try:
    api_response=core_api.create_namespaced_persistent_volume_claim(namespace='devops',body=bodypvc)
    print(api_response)
except Exception as e:
    print("Exception when calling CoreV1Api->create_namespaced_persistent_volume_claim: %s\n" % e)

```
