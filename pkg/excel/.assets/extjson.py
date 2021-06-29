
#%%
import json
import random
MaxLine = 500

jmap = json.load(open("list.json", "r", encoding="utf8"))
outFp = open("list-%d.json"%(MaxLine) , "w", encoding="utf8")
outList = [jmap[0]]

for i in range(MaxLine):
    outList.append(jmap[random.randint(1, len(jmap) - 1)])

json.dump(outList, outFp)

# %%
