# golang-challenge

###### Consumiendo la API de Marvel

#### Usabilidad y flujo de la aplicación
- La aplicación permite indicar tokens propios o usar los predeterminados
- Se pide seleccionar una de dos opciones.
- 1.- Buscar por nombre: Por dentro este funciona como `nameStartsWith` para búsquedas más eficientes
- 2.- Listar los personajes: Se extrae una lista con los primeros 20 datos devueltos
- Si no se elige ninguna y omite con un salto de línea se inicia la tarea por defecto (se puede cambiar los parámetros por defecto ajustando la función `searchParameters()`)
- Al terminar alguna petición realizada se pregunta si desea seguir buscando, en caso de ser sí, el proceso continuará usando las tokens anteriormente declaradas.

#### Recomendaciones
- Si la aplicación falla leyendo tus entradas de teclado puedes cambiar el global  `inputrefactor` para quitar algún `final de línea` de tu consola
- Algunos nombres de super heroes vienen con guión. Ejemplo: `spider-man`
- Si deseas cambiar las tokens del source usa las `flags: {apiKey} y {secret}`
