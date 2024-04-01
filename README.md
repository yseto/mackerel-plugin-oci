# yseto/mackerel-plugin-oci

it is alpha quality.

## usage

### need credentials

need `~opc/.oci/...`

## sample of mackerel-agent.conf

```
[plugin.metrics.ocimds]
user = "opc"
command = ["/path/to/mackerel-plugin-oci-mds", "-compartmentId", "ocid1.compartment.oc1..xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "-resourceId", "ocid1.mysqldbsystem.oc1.ap-tokyo-1.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"]

[plugin.metrics.flb]
user = "opc"
command = ["/path/to/mackerel-plugin-oci-flb", "-compartmentId", "ocid1.compartment.oc1..xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "-resourceId", "ocid1.loadbalancer.oc1.ap-tokyo-1.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"]

[plugin.metrics.nlb]
user = "opc"
command = ["/path/to/mackerel-plugin-oci-nlb", "-compartmentId", "ocid1.compartment.oc1..xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "-resourceId", "ocid1.networkloadbalancer.oc1.ap-tokyo-1.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "-resourceName", "xxx-nlb"]
```

