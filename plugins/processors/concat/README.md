# Concat Processor Plugin

The `concat` plugin concatenates multiple tag values into a new tag, optionally using a separator to combine the tag values into a single string.

### Configuration:

```toml
[[processors.concat]]
  namepass = ["interface"]

  ## Tag concatenations defined in a separate sub-table
  [[processors.concat.tags]]
    ## Tags that should be concatenated
    keys = [ "host", "if_name" ]
    ## Separator used to join the values together
    separator = ":"
    ## New tag key
    result_key = "global_if_name"
```

### Tags:

No tags are applied by this processor.

### Example Output:
```
interface,host=my.host,if_name=eth1,global_if_name=my.host:eth1 octects_in=51684535,octets_out=5457728 1519652321000000000
```
