# OS-diff
Openstack / Openshift diff tool

This tool collects Openstack/Openshift service configurations,
compares configuration files, makes a diff and creates a report to the user
in order to provide informations and warnings after a migration from
Openstack to Openstack on Openshift migration.

### Usage

#### Pull configuration step

Before running the Pull command you need to configure the ssh access to your environements (Openstack and OCP).
Edit the ssh.config provided with this project and make sure you can ssh on your hosts with the command:

```
ssh -F ssh.config crc
ssh -F ssh.config standalone
```

Also you need to provide the full path of your ssh.config in the ansible.cfg file, example:

```
ssh_args = -F /home/foo/os-diff/ssh.config
```

When everything is setup correctly you can tweak the ansible vars for each services you want to analyze:

```
  ▾ roles/
    ▾ collect_config/
      ▾ vars/
        main.yml
```

You can add your own service according to the following:

```
  # Service name
  keystone:
    # Bool to enable/disable a service (not implemented yet)
    enable: true
    # Pod name, in both OCP and podman context.
    # It could be strict match with strict_pod_name_match set to true
    # or by default it will just grep the podman and work with all the pods
    # which matched with pod_name.
    pod_name: keystone
    # Path of the config files you want to analyze.
    # It could be whatever path you want:
    # /etc/<service_name> or /etc or /usr/share/<something> or even /
    # @TODO: need to implement loop over path to support multiple paths such as:
    # - /etc
    # - /usr/share
    path: /etc/keystone
    # In podman context, when you want to pull specific files:
    # You need to set pull_items to true
    name:
      - keystone.conf
      - logging.conf
```

An Ansible hosts file is provided at the root of this repository and the
ansible.cfg.
You might want to edit the hosts file to stick to your environment.
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

Once everything is correctly setup you can start to pull configuration:


```
# install dependencies
make install
# build os-diff
make build
# run pull configuration for TripleO standalone:
./os-diff pull --cloud_engine=podman --inventory=$PWD/hosts
# run pull configuration for OCP:
./os-diff pull --cloud_engine=ocp --inventory=$PWD/hosts

# You can also use the playbooks directly:
ansible-playbook -i hosts playbooks/collect_ocp_config.yaml
```

#### Compare configuration files steps

Once you have collected all the data per services you need, you can start to run comparison between
your two source directories.
A results file is written at the root of this project `results.log` and a *.diff file is created for each
file where a difference has been detected

```
/tmp/collect_crc_configs/nova/nova-api-0/etc/nova/nova.conf.diff

# with this kind of content:
Source file path: /tmp/collect_crc_configs/nova/nova-api-0/etc/nova/nova.conf, difference with: /tmp/collect_crc_configs/nova/nova-cell0-conductor-0/etc/nova/nova.conf
[DEFAULT]
-transport_url=rabbit://default_user_pVPGFkYMWTdSarUSog9:Rg59ofmjeDWg24v8ZeGW-1PblH1LJDQ1@rabbitmq.openstack.svc:5672
[api]
-auth_strategy=keystone
```

The log INFO/WARN and ERROR will be print to the console as well so you can have colored info regarding the current file processing.
Run the compare command:

```
./os-diff compare --origin=/tmp/collect_tripleo_configs --destination=/tmp/collect_crc_configs

```

### Asciinema demo

https://asciinema.org/a/5YnpA5t7uZKn2H2jU5LSNHNQx


### TODO

* Add option to compare the config files directly from files to pods
* Improve reporting (console, debug and log file with general report)
* Improve diff output for json and yaml
* Improve Makefile entry with for example: make compare
* Add a skip list (skip /etc/keystone/fernet-keys )
* Add interactive and edit mode to ask for editing the config for the user
  when a difference has been found

