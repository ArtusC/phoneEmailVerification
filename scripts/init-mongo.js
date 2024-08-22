print("[Creating collection]: get phoneDb Database ")
db = db.getSiblingDB('phoneDb')

print("[Creating collection]: phone-collection ")
db.createCollection("phone-collection")