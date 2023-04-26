# OS-diff
Openstack / Openshift diff tool

This tool collects Openstack/Openshift service configurations,
compares configuration files, makes a diff and creates a report to the user
in order to provide informations and warnings after a migration from
Openstack to Openstack on Openshift migration.

### Usage

An Ansible hosts file is provided at the root of this repository and the
ansible.cfg.
Those file are required for collecting the configuration files from
the pods or the containers (OCP and Podman).

```
  ▾ playbooks/
      collect_ocp_config.yaml
      collect_podman_config.yaml
```

Those playbooks can call with the Go binary or directly with Ansible.
It call one Ansible role:

```
  ▾ roles/
    ▾ collect_config/
      ▾ tasks/
        collect_ocp.yml
        collect_podman.yml
        main.yml
```

A complete CLI provided with the Go binary will be available soon.
For now you can execute those playbooks with this command:

```
os-diff pull --cloud_engine=ocp --inventory=$PWD/hosts

```

And execute the comparision of your configuration files with (edit the go files
to provide a valid path for your datas):

```
os-diff compare --origin=tests/podman/keystone.conf --destination=tests/ocp/keystone.conf --output=output.txt

```

### TODO

* Add option to compare the config files directly from files to pods
* Compare directories
* Handle other config type: yaml, json.
