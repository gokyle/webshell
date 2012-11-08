/*
   Copyright (c) 2012 Kyle Isom <kyle@tyrfingr.is>

   Permission to use, copy, modify, and distribute this software for any
   purpose with or without fee is hereby granted, provided that the 
   above copyright notice and this permission notice appear in all 
   copies.

   THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL 
   WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED 
   WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE 
   AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL
   DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA
   OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER 
   TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR 
   PERFORMANCE OF THIS SOFTWARE.
*/

/*
   webshell is an outline for a webserver to save the 10 minutes it would
   take to type all this up. it is designed to minimise the time it takes
   me to get started on projects and to enable one-off temporary / single
   purpose webservers.

   The environment variables SSL_KEY and SSL_CERT should point to the SSL
   key and cert if the server is to listen for TLS connections.

   The environment variables SERVER_ADDR and SERVER_PORT should contain
   the listening address and port; at a minimum, the port is required. It
   is acceptable to specify the port in SERVER_ADDR; if so, SERVER_PORT
   must be left blank.
*/
package webshell
