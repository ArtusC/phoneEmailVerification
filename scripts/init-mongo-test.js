print("[Creating collection]: get phoneDb Database ")
db = db.getSiblingDB('phoneDb_test')

print("[Creating collection]: phone-collection ")
db.createCollection("phone-collection-test")