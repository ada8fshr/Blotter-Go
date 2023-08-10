import utils

old_data = utils.query("select * from tags;")

new_data = list(map(lambda x: {
    "short": x["english"].lower(),
    "name": x["chinese"],
    "color": "",
    "icon": x["img"],
}, old_data))

for d in new_data:
    print(d)

document = utils.mydb["tags"]
document.delete_many({})
ids = document.insert_many(new_data)
