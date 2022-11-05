package main

func Example() {
	main()
	// Output:
	// clients.csv:
	// client_id|client_name|client_age|addr1.street|addr1.city|Address2.street|Address2.city|employed
	// 12|John|21|Street 1|City1|Street 2|City2|20220411
	// 13|Fred||"Main ""Street"" 1"|City1|Main Street 2|City2|20220411
	// 14|James|32|Center Street 1|City1|Center Street 2|City2|20220411
	// 15|Danny||State Street 1|City1|State Street 2|City2|20220411
	//
	// clients:
	// 12:John Adress1:Street 1, City1 Address2:Street 2, City2 Employed: 20220411
	// 13:Fred Adress1:Main "Street" 1, City1 Address2:Main Street 2, City2 Employed: 20220411
	// 14:James Adress1:Center Street 1, City1 Address2:Center Street 2, City2 Employed: 20220411
	// 15:Danny Adress1:State Street 1, City1 Address2:State Street 2, City2 Employed: 20220411
}
