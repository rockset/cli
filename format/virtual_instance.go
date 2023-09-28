package format

var VirtualInstanceDefaultSelector = DefaultSelector{
	Normal: "Name:.name,Description:.description,State:.state,Default VI:.default_vi,Current Size:.current_size",
	Wide:   "Name:.name,ID:.id,description:.description,State:.state,Default VI:.default_vi,Current Size:.current_size,Desired Size:.desired_size",
}
