blocked = []
x = """
$ORIGIN test.foo.
$TTL 1750
@	IN	SOA	master.test.foo (
	2019040502	; serial
	86400		; refresh
	7200		; retry
	3600000	; expire
	1750		; minimum
)
	IN	NS	ns.test.foo.
www	IN	A	93.184.216.34
"""
def readFile(fobj):
    rtn = []
    lines = fobj.readlines()                                                                                                                                                   
    for line in lines:                                                                                                                                                        
        line = line.strip()                                                                                                                                                   
        if line.startswith("zone"):                                                                                                                                           
            tokens = line.split(" ")                                                                                                                                          
            token = tokens[1][1:-1]                                                                                                                                           
            if token not in rtn:
                rtn.append(token)
        elif line.startswith("0.0.0.0"):
            tokens = line.split(" ")
            token = tokens[1]
            if token not in rtn:
                rtn.append(token)
        # print(len(rtn))
    fobj.close()
    return rtn 
           
fb0 = open("/etc/coredns/hosts.blocked", "r")
domains = readFile(fb0)
fb1 = open("/blocked/malware.zone", "r")
ds = readFile(fb1)
for d in ds:
    if d not in domains:
       domains.append(d)
fb2 = open("/blocked/sb-hosts/data/StevenBlack/hosts", "r")
ds = readFile(fb2)
for d in ds:
    if d not in domains:
       domains.append(d) 
print(len(domains))

