- Wenn man einen DetailedError wrapt (also in einen neuen Detailed error), sollte der neue
  DetailedError auch einen StackTrace generieren oder einfach den von dem gewrapten
  DetailedError übernehmen?

- Was soll der default bei einem Detailed error sein für die `Error()` Methode? Alle details oder
  nur die eigentliche Error message sodass man für alle details eine eigene Methode hat?
  In dem Fall könnte man eine Funktion `Details() string` einbauen, welche, wenn ein DetailedError
  übergeben wird, alle Details formatiert zurückgibt und bei allen anderen Errors nur die message
  zurückgibt.