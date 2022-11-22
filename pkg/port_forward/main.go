package port_forward

import "regexp"

var PortMappingExp = regexp.MustCompile("^(?P<local>[1-9][0-9]*)?(:(?P<remote>[1-9][0-9]*))?$")
