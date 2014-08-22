# myjsonip.com in Go

This is a site that allows you to quickly get information relative to your hosts. IP address and agent are the only things
implemented so far.


## How to use the site

You are able to query the site for your WAN IP address or agent and it will return it in json (default), yaml or xml depending on the route.

| HTTP Method | URI            | What is returned?  |
|:------------|:---------------|:-------------------|
| GET         | /:format       | IP Address         |
| GET         | /ip/:format    | IP Address         |
| GET         | /agent/:format | Agent              |
| GET         | /all/:format   | IP Address + Agent |

### Examples:

http://myjsonip.com/ip or http://myjsonip.com/ip/json

```json
{"ip":"172.22.52.123"}
```

http://myjsonip.com/ip/yaml:

```yaml
---
ip: 172.22.52.123
```

http://myjsonip.com/ip/xml:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<hash>
  <ip>172.22.52.123</ip>
</hash>
```

## I have _x_ interfaces on my server with NATs, how can I leverage this site to get each interface WAN IP?

I use the following script in production in order to quickly look up all interfaces or a single interface WAN
IP. A few examples of how to use the script:

How to grab the WAN IP for ALL interfaces: `$ find_wan_ip.rb`

How to grab the WAN IP for a single interface: `$ find_wan_ip.rb 192.168.2.5`

How to grab the WAN IP for multiple interfaces: `$ find_wan_ip.rb 192.168.2.5 192.168.2.6`

    require 'net/http'
    require 'json'
    require 'resolv'
    require 'socket'

    # Grab all interfaces that have valid IPv4 addresses
    interfaces = Socket.getifaddrs.map do |i|
      next if i.addr.ipv4_loopback? || i.addr.ipv6_loopback?
      next unless i.addr.ipv4?
      x = [ i.name, i.addr.ip_address ]
      x
    end.compact

    # Build hash for looking up interface names by IP
    ip_to_int = interfaces.inject({}) do |memo,x|
      name   = x[0]
      lan_ip = x[1]
      memo[lan_ip] = name
      memo
    end

    # Parse the command line to see if any IP addresses where passed.
    ip_addresses = ARGV.each.map do |x|
      if x =~ Resolv::IPv4::Regex && ip_to_int.keys.include?(x)
        x
      else
        next
      end
    end.compact

    # If no IP addresses where passed, grab the WAN IP of all interfaces
    ip_addresses = interfaces.map { |x| x[1] } if ip_addresses.empty?

    # Setup our http connection. This is used to discover the public/WAN IP.
    uri  = URI('http://myjsonip.com/ip')
    http = Net::HTTP.new(uri.host, uri.port)
    http.open_timeout = 2
    http.read_timeout = 2

    ip_addresses.each do |ip|
      failed_request = false
      no_ip = false
      lan_ip = ip
      interface_name = ip_to_int[lan_ip]

      # This makes sure we use this interface when making outbound connections.
      http.local_host = lan_ip

      begin
        res = http.get('/')
      rescue Net::OpenTimeout, Errno::EHOSTUNREACH, Errno::ECONNREFUSED => e
        failed_request = true
        no_ip = true
      end

      wan_ip   = JSON.load(res.body)['ip'] unless failed_request
      # If the WAN IP is 50.56.18.108 it means that no NAT is setup and it's defaulting to our firewall's IP.
      wan_ip = nil if no_ip
      puts "%s:%s:%s" % [interface_name, lan_ip, wan_ip]
    end

## Contributing

1. Fork it ( https://github.com/jasonwbarnett/myjsonip.com-go/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
