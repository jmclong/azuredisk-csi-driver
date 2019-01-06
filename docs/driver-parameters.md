### `disk.csi.azure.com` driver parameters
 > storage class `disk.csi.azure.com` parameters are compatible with built-in [azuredisk](https://kubernetes.io/docs/concepts/storage/volumes/#azuredisk) plugin

 - Dynamic Provisioning
  > get a quick example [here](../deploy/example/storageclass-azuredisk-csi.yaml)

Name | Meaning | Available Value | Mandatory | Default value
--- | --- | --- | --- | ---
skuName | azure disk storage account type (alias: `storageAccountType`)| `Standard_LRS`, `Premium_LRS`, `StandardSSD_LRS`, `UltraSSD_LRS` | No | `Standard_LRS`
kind | managed or unmanaged(blob based) disk | `managed` (`dedicated`, `shared` are deprecated since it's using unmanaged disk) | No | `managed`
fsType | File System Type | `ext4`, `ext3`, `xfs` | No | `ext4`
cachingMode | azure data disk host cache setting | No | `None`, `ReadOnly`, `ReadWrite` | `ReadOnly`
storageAccount | specify the storage account name in which azure disk will be created | STORAGE_ACCOUNT_NAME | No | if empty, driver will find a suitable storage account that matches `skuName` in the same resource group as current k8s cluster
location | specify the Azure location in which azure disk will be created | `eastus`, `westus`, etc. | No | if empty, driver will use the same location name as current k8s cluster
resourceGroup | specify the resource group in which azure disk will be created | existing resource group name | No | if empty, driver will use the same resource group name as current k8s cluster
DiskIOPSReadWrite | [UltraSSD disk](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/disks-ultra-ssd) IOPS Capability | 100~160000 | No | `500`
DiskMBpsReadWrite | [UltraSSD disk](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/disks-ultra-ssd) Throughput Capability | 1~2000 | No | `100`

 - Static Provisioning(use existing azure disk)
 
Name | Meaning | Available Value | Mandatory | Default value
--- | --- | --- | --- | ---
volumeHandle| azure disk URI | /subscriptions/{sub-id}/resourcegroups/{group-name}/providers/microsoft.compute/disks/{disk-id} | Yes | N/A
volumeAttributes.fsType | File System Type | `ext4`, `ext3`, `xfs` | No | `ext4`
volumeAttributes.cachingMode | azure data disk host cache setting | No | `None`, `ReadOnly`, `ReadWrite` | `ReadOnly`

  > get a quick example [here](../deploy/example/pv-azuredisk-csi.yaml) 