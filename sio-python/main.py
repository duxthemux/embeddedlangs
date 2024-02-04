import sys, json;
data = json.load(sys.stdin)

data["field3"]="nono"

print(json.dumps(data))

sys.stdout.write("{}")