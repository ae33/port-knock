# port-knock
A tool to easily perform arbitrary port-knocking.

## Setup

#### Client-side
First, fork and clone this repo.

Next, you will need to define a config file. You can find examples in 
`config/examples`. Your final config file should be located at 
`/etc/port-knock/port-knock.yml`. 

You do have the option to use a different location for the config file via the 
`-config-path` flag to `main.go`, but that requires you to either invoke 
`main.go` with `go run`, or to edit the source code, and build the binary with 
`go build`. In the latter scenario, you'll need to make sure you have a proper 
Golang environment set up.
 
After finalizing your config, take a look at `scripts/knock_template.sh`. Make 
sure you have execute permissions on the script (e.g. 
`chmod 744 scripts/knock_template.sh`). Then, edit the final line in the script 
to have your ssh information.

#### Server-side
You will need to make sure your ssh daemon is running. This varies from OS to 
OS, but in Linux distributions, it's likely that you have a config in 
`/etc/ssh`. Make sure you edit your server config, and not your client config.
In Arch Linux, for example, you would need to edit `/etc/ssh/sshd_config`. Then,
you would start the service that references that config by executing:
```bash
systemctl enable sshd
systemctl start sshd
```

Next, you'll need to make sure your firewall will allow all of your port knocks 
through and will allow the SSH connection.

First, try to get an ssh connection without port knocking by running the 
following commands as root:
```bash
# Allow all outbound traffic from your server.
iptables -A OUTPUT -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT

# Allow inbound TCP traffic to your SSH port.
iptables -A INPUT -p tcp --dport <ssh-port> -m state --state NEW -j ACCEPT
```

Then, try to ssh in to your SSH server from some other computer. Make sure you 
have port-forwarding properly set up in your router's settings.

Once you've gotten that working, execute the following to set up a three-knock 
port knocking scheme:
```bash
# Define the chains for the desired port-knocking sequence.
iptables -N STATE0
iptables -A STATE0 -p udp --dport <first-knock-port> -m recent --name KNOCK1 --set -j DROP
iptables -A STATE0 -j DROP

iptables -N STATE1
iptables -A STATE1 -m recent --name KNOCK1 --remove
iptables -A STATE1 -p udp --dport <second-knock-port> -m recent --name KNOCK2 --set -j DROP
iptables -A STATE1 -j STATE0

iptables -N STATE2
iptables -A STATE2 -m recent --name KNOCK2 --remove
iptables -A STATE2 -p udp --dport <third-knock-port> -m recent --name KNOCK3 --set -j DROP
iptables -A STATE2 -j STATE0

iptables -N STATE3
iptables -A STATE3 -m recent --name KNOCK3 --remove
iptables -A STATE3 -p tcp --dport <ssh-port> -j ACCEPT
iptables -A STATE3 -j STATE0

# Add the port-knocking chains to the main INPUT chain.
iptables -A INPUT -m recent --name KNOCK3 --rcheck -j STATE3
iptables -A INPUT -m recent --name KNOCK2 --rcheck -j STATE2
iptables -A INPUT -m recent --name KNOCK1 --rcheck -j STATE1
iptables -A INPUT -j STATE0
```

Make sure to save these `iptables` rules, so they persist. In Arch Linux, for
example, you would execute: 
```bash
iptables-save > /etc/iptables/iptables.rules
```

## Usage
Execute your edited `scripts/knock_template.sh`, and voila! You're in.

You could even consider adding the script to your PATH, for easy access. A 
suitable place would be in `/usr/bin`. If you do this, you could use a much
simpler script than `scripts/knock_template.sh` and simply have a call to your
binary, followed by your ssh command. If you alias this script, you could have a
one word command to connect to your SSH server.

## Troubleshooting
If your script is hanging, it is possible that your router isn't forwarding all
your ports, properly. Make sure that you have port forwarding set up for all the
port-knocking ports under the UDP protocol, and the TCP protocol for the SSH 
port.

It is also possible that your SSH server's firewall isn't set up properly. Check 
your active rules with `iptables -S` or `iptables -L`, and make sure you aren't 
blocking any of your desired inbound or outbound traffic.

A basic `iptables` cheat sheet:

| Flag | Description
|------|-----------------------------------------------------------------------------------------|
| -P   | Defines a new default policy.                                                           |
|  -N  | Defines a new chain.                                                                    |
| -A   | Adds to an existing chain.                                                              |
| -m   | Match the given rule against the given traffic.                                          |
| -j   | Jump to a target. There are final targets of ACCEPT, DROP, and RETURN. |
