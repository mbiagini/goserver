# Dev

## Investigaciones:

1. Agregar validators.
   1. DONE. Se usa https://github.com/go-playground/validator dentro del package gsvalidation.
2. Ver Oauth
3. Ver throttling [https://github.com/martini-contrib/throttle]
4. Investigar cómo funciona la carpeta vendor (go mod vendor).
   1. DONE. Al momento de compilar el código, el compilador de GO va a buscar primero las dependencias en vendor y, si no encuentra una, recíen ahí la va a buscar en $GOPATH y $GOROOT.
5. Ver qué es Godoc.
   1. DONE. Sirve para generar swagger y documentación estática.
6. Ver la estructura cmd, pkg, internal para estructurar el proyecto.
7. Ver cómo exponer servicios como dependencias (https://medium.com/@ott.kristian/how-i-structure-services-in-go-19147ad0e6bd).

## TODO:

1. Usar Swagg para generar swagger y servirlo con swagger-ui (ver que se haga todo automáticamente, sin tener que generar el .json y luego servirlo).
2. Agregar métricas de prometheus con https://prometheus.io/docs/guides/go-application/.
   1. Agregar métricas básicas sin que aparezcan los requests en el http_logger.
   2. Agregar un Transport que registre métricas de forma automática (contando invocaciones hacia cada backend exitosas y de error).
   3. Agregar un middleware (Handler) que registre métricas automáticas de invocaciones exitosas y erróneas al servicio.
3. Agregar swagger-ui y ejecución de swag init en el run.sh.
4. Separar la carpeta utils como dependencia.
5. Agregar validación para parámetros simples en gsvalidation.
6. Exponer el swagger interno a través de una url /swagger.json (ponele). Si es posible, con swagger-ui.
7. Hacer una librería para consultas sql que tenga validaciones útiles.
8. Ver cómo evitar el GetClient by key en cada operación de un client.
   1. Ver de llevar el type ClientKey a gsclient. El único problema por ahora es que para crear un nuevo client, se usa la config. Entonces, al parsear el json con el string "MockClient", debería traducirse en la constante...
