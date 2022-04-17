echo "All users"
curl "localhost:10000/users"

echo "\n Add first user Petrov Petr Petrovich"
curl --request POST "localhost:10000/users" --data '{"username" : "Petrov Petr Petrovich"}'

echo "\n Add second user Ivanov Ivan Asetrovich"
curl --request POST "localhost:10000/users" --data '{"username" : "Ivanov Ivan Asetrovich"}'

echo "\n Add third user Alex Alexov Alexovich"
curl --request POST "localhost:10000/users" --data '{"username" : "Alex Alexov Alexovich"}'

echo "\n Add fourth user Vlad Kek Alexovich"
curl --request POST "localhost:10000/users" --data '{"username" : "Vlad Kek Alexovich"}'

echo "\n Get list of users"
curl "localhost:10000/users"

echo "\n Get info about first user"
curl "localhost:10000/users/1"

echo "\n Get info about non-existent user"
curl "localhost:10000/users/101"

echo "\n Edit info about 2 user to Some Some Somisch"
curl --request PUT "localhost:10000/users/2" --data '{"username" : "Some Some Somisch"}'

echo "\n Edit info about non-existent user to Some Some Somisch"
curl --request PUT "localhost:10000/users/1000" --data '{"username" : "Some Some Somisch"}'

echo "\n Delete info about 1 user"
curl --request DELETE "localhost:10000/users/1"

echo "\n Delete info non-existent 1 user"
curl --request DELETE "localhost:10000/users/100"

echo "\n Get list of users"
curl "localhost:10000/users"