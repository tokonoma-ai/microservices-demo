function get_results(result) {
    print(tojson(result));
}

function insert_address(object) {
    print(db.addresses.insert(object));
}

insert_address({
    "_id": ObjectId("57a98d98e4b00679b4a830ad"),
    "number": "246",
    "street": "Whitelees Road",
    "city": "Glasgow",
    "postcode": "G67 3DL",
    "country": "United Kingdom"
});
insert_address({
    "_id": ObjectId("57a98d98e4b00679b4a830b0"),
    "number": "246",
    "street": "Whitelees Road",
    "city": "Glasgow",
    "postcode": "G67 3DL",
    "country": "United Kingdom"
});
insert_address({
    "_id": ObjectId("57a98d98e4b00679b4a830b3"),
    "number": "4",
    "street": "Maes-Y-Deri",
    "city": "Aberdare",
    "postcode": "CF44 6TF",
    "country": "United Kingdom"
});
insert_address({
    "_id": ObjectId("57a98ddce4b00679b4a830d1"),
    "number": "3",
    "street": "Rochester Ave",
    "city": "London",
    "country": "UK"
});
insert_address({
    "_id": ObjectId("57a98d98e4b00679b4a830c0"),
    "number": "12",
    "street": "Princes Street",
    "city": "Edinburgh",
    "postcode": "EH2 2AN",
    "country": "United Kingdom"
});
insert_address({
    "_id": ObjectId("57a98d98e4b00679b4a830c3"),
    "number": "87",
    "street": "Deansgate",
    "city": "Manchester",
    "postcode": "M3 2BW",
    "country": "United Kingdom"
});

print("________ADDRESS DATA_______");
db.addresses.find().forEach(get_results);
print("______END ADDRESS DATA_____");
