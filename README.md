# Deal Scraper

Acest repository conține codul pentru backend-ul Deal Scaper
Deal Scraper este o aplicație care îți scanează galeria după imagini cu bonuri fiscale, 
identifică produsele cumpărate și caută cele mai bune oferte pentru produsele pe care le-ai cumpărat.

Pentru realizarea acestui scop. aplicația se folosește de tehnici de învățare automată și web scraping.

Backend-ul este construit folosind Go datorită performaței crescute ale limbajului și suportul foarte bun pentru programare concurentă 
prin intermediu gorutinelor. 

Aplicația este construită folosind o arhitectură bazată pe microservicii. Aceste microservcii comunică între ele atât sincron
prin intermediul cererilor HTTP REST, cât și asincron prin intermediul protocolului AMQP. Pentru implementarea serverelor HTTP, Deal Scraper
se folosește de framework-ul [Gin](https://github.com/gin-gonic/gin), ales pentru minimalismul și ușurința sa de implementare, 
avâmd totuși funcționalitățile necesare pentru Deal Scraper și performanță foarte bună. Pentru comunicarea asincronă este folosit RabbitMQ, 
o implementare foarte populară pentru protocolul AMQP. De asemenea, pentru persistența datelor am folosit baze de date relaționare prin MySQL.

## Microservicii

Aplicația se folosește de mai multe microservicii pentru realizarea scopului său, și anume identificarea produselor cumpărate de utilizatori
din imaginile acestora și găsirea de oferte mai bune pentru acele produse. Microserviciile folosite sunt următoarele:

- API Gateway -> Abstractizează toate microserviciile pentru aplicația client în spatele unui singur API REST
- Auth -> Se ocupă de autentificarea și datele utilizatorilor
- Store Metadata -> Expune un API REST pentru accesarea unei baze de date care conține informații despre magazinele din România
- OCR -> Folosind Google Cloud Vission extrage textul din imaginile utilizatorului pentru a obține lista de produse achiziționate
- Search -> Folosind Google Search, acest microserviciu caută oferte la mai multe magazine pentru produsele extrase de OCR și trimite
URL-urile către aceste oferte mai departe în pipeline
- Scheduler -> Folosind un CRON programează când și care oferte trebuiesc crawluite
- Crawler -> Conține web crawlers pentru mai multe magazine și este folosit pentru extragere ofertelor
- Product Processing -> Procesează ofertele extrase, realizează un grad de legături între produse și oferte pe care-l salvează în baza de date
 și expune aceste oferte către utilizatori
- Notifications -> Se folosește de Firebase Cloud Messaging pentru a trimite notificări push către utilizatori când noi oferte pentru
produsele pe care aceștia le-au cumpărat sunt disponibile

## Diagramă flow microservicii

[Diagramă.pdf](https://github.com/adiRandom/grocever-backend/files/11725700/Licenta_diagrama.pdf)
