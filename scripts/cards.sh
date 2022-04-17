echo "All cards"
curl "localhost:10000/cards"

echo "\n Add first card with userId=2"
curl --request POST "localhost:10000/cards" --data '{"balance" : 1000, "userId" : 2}'

echo "\n Add second card with userId=3"
curl --request POST "localhost:10000/cards" --data '{"balance" : 2000, "userId" : 3}'

echo "\n Add third card with userId=4"
curl --request POST "localhost:10000/cards" --data '{"balance" : 3000, "userId" : 4}'

echo "\n Get list of cards"
curl "localhost:10000/cards"

echo "\n Get info about first card"
curl "localhost:10000/cards/1"

echo "\n Get info about non-existent card"
curl "localhost:10000/cards/101"

echo "\n Edit info about 1 card (balance to 5000)"
curl --request PUT "localhost:10000/cards/1" --data '{"balance" : 5000}'

echo "\n Edit info about non-existent card"
curl --request PUT "localhost:10000/cards/101" --data '{"balance" : 5000}'

echo "\n Get list of cards"
curl "localhost:10000/cards"

echo "\n Delete info about first card"
curl --request DELETE "localhost:10000/cards/1"

echo "\n Delete info about non-existent card"
curl --request DELETE "localhost:10000/cards/1"

echo "\n Get list of cards"
curl "localhost:10000/cards"

echo "\n Refill second card"
curl --request POST "localhost:10000/cards/2" --data '{"AddBalance" : 100000}'

echo "\n Refill non-existent card"
curl --request POST "localhost:10000/cards/200" --data '{"AddBalance" : 100000}'

echo "\n Transfer from one to other card (from second card to third card, 100 counts"
curl --request POST "localhost:10000/cards/transfer" --data '{"CardFrom" : 2, "CardTo": 3, "AddBalance" : 100}'

echo "\n Negative case of transfer from one to other card"
curl --request POST "localhost:10000/cards/transfer" --data '{"CardFrom" : 3, "CardTo": 4, "AddBalance" : 100000}'

echo "\n Get list of cards"
curl "localhost:10000/cards"