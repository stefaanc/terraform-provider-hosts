###
### myhost1
###

output "myhost1_name" {
    value = "${data.hosts_record.myhost1.name}"
}

output "myhost1_record_id" {
    value = "${data.hosts_record.myhost1.record_id}"
}

output "myhost1_address" {
    value = "${data.hosts_record.myhost1.address}"
}

output "myhost1_names" {
    value = "${data.hosts_record.myhost1.names}"
}

output "myhost1_comment" {
    value = "${data.hosts_record.myhost1.comment}"
}

###
### myhost111
###

output "myhost111_record_id" {
    value = "${hosts_record.myhost111.record_id}"
}

output "myhost111_address" {
    value = "${hosts_record.myhost111.address}"
}

output "myhost111_names" {
    value = "${hosts_record.myhost111.names}"
}

output "myhost111_comment" {
    value = "${hosts_record.myhost111.comment}"
}

output "myhost111_notes" {
    value = "${hosts_record.myhost111.notes}"
}
