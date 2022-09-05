# Test task for links searcher implementation

Implement utility for weblinks search and validation (the link to site must be working). 
Program reads lines from stdin and finds links in it, then checks if they are working and at the end of program work writes to stdout all the valid links.
Lines must be processed concurrently, but not more than N at a time (N - configured value), and links search must start immediately after its reading
If input or output files aren't provided, use standart input and output