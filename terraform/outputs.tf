output "myhost1_address" {
    value = "${data.hosts_record.myhost1.address}"
}

output "myhost1_name" {
    value = "${data.hosts_record.myhost1.name}"
}