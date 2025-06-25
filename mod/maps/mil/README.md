# maps



## routing

Problems:

1. Its using gio, which will NOT cross compile so a PITA
Cogent or Ebiten gugui

2. The SQLITE uses weird bindings. Try to use the wasm version
https://github.com/ncruces/go-sqlite3

CSV 
- sqlite3_auto_extension((void*) sqlite3_csv_init);
- https://github.com/ncruces/go-sqlite3/blob/main/ext/csv/csv.go has it and runs on any device without any CGO

! Ebiten works as wasm or native, which means way less screwing around with html and css. Gemini also likes it



https://git.sr.ht/~mil/

m@milesalan.com

https://mr.lrdu.org

3 primary components:

Mobsql (GTFS-to-SQLite ETL pipeline: CLI + Go Library)
Mobroute (Core routing algorithm via CSA: CLI + Go Library)
Transito (Graphical mobile application for Android & Linux)


https://git.sr.ht/~mil/mobsql is the data layer using sqlite.
https://git.sr.ht/~mil/mobsql/tree

https://sr.ht/~mil/mobroute/ is the logic
https://git.sr.ht/~mil/mobroute/tree

https://sr.ht/~mil/transito/ is the GIO gui
https://git.sr.ht/~mil/transito/tree
 

## readeck



https://codeberg.org/readeck/readeck

```sh
docker run --rm -ti -p 8000:8000 -v readeck-data:/readeck codeberg.org/readeck/readeck:latest
```  


