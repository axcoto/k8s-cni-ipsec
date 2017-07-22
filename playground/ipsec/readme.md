# Note

## Client talk to each other via virtual ip

1. Strongswan server has to allow forward package(it's 0 by default)

  echo 1 > /proc/sys/net/ipv4/ip_forward

2. Add virtual ip into subnet, this will make strongswan insert necessary route into 220 table

   Something like this on client, when 10.100.255.0 is virtual ip pool in this case

   rightsubnet=10.0.0.0/16,10.100.255.0/28

   Respectively, on the server:

   rightsubnet=10.0.0.0/16,10.100.255.0/28

   With this setup, strongswan create this route on client:

   ```
   10.0.0.0/16 via 10.0.2.2 dev enp0s3  proto static  src 10.100.255.2
   10.100.255.0/28 via 10.0.2.2 dev enp0s3  proto static  src 10.100.255.2
   ```

## To use same password for everything(not secure of course). just getting done

%any : PSK "secret"
