# stripe

## code

api: https://github.com/stripe/stripe-go
api docs: https://docs.stripe.com/api?lang=go 


github.com/stripe/stripe-go/v82 seems to ve what the v2 code uses !

---

cli: https://github.com/stripe/stripe-cli
cli docs: https://docs.stripe.com/cli

## tech

### server

My Sever takes calls from the GUI, talks to Stripe, and tracks things using events and updates the Server DB and the Golang GUI.

### client


GUI needs to access users payment info from golang.

https://github.com/progrium/darwinkit

https://github.com/tractordev/toolkit-go/tree/main/desktop



used by: https://github.com/tractordev/toolkit-go/network/dependents?package_id=UGFja2FnZS00MTAwOTY1ODAw

https://github.com/lonelyeel/go-vscode/blob/master/_example/main.go has OPFS !

https://github.com/tractordev/toolkit-go/tree/main/engine/fs


## poc

Musa Ramli
mramli-at-sales.stripe.com@growth.stripe.com

## docs

### locations

https://docs.stripe.com/tax/supported-countries

For example, Australia is https://docs.stripe.com/tax/supported-countries/asia-pacific/australia

For example, in Australia, they support Digital and physical goods and Remote sales. 
Also for all of the EU and US you have full support. 

I need to do "Drop shipping" into these countries for the physical goods. That's when the physical goods are delivered from the Factory ( in China in my case  ) direct to the Customers address



## examples




