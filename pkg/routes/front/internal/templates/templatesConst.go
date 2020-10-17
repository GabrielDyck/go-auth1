package templates

const ErrorTemplate string = `
<html>
<body>
<h1> Error ocurred </h1>
<h2> {{.Reason}}
</body>
</html>`


const RedirectAuthenticationError= `
<html>
<head> 
    <script src="//ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>

<script>

        function getCookie(name) {
            // Split cookie string and get all individual name=value pairs in an array
            var cookieArr = document.cookie.split(";");

            // Loop through the array elements
            for(var i = 0; i < cookieArr.length; i++) {
                var cookiePair = cookieArr[i].split("=");

                /* Removing whitespace at the beginning of the cookie name
                and compare it with the given string */
                if(name == cookiePair[0].trim()) {
                    // Decode the cookie value and return
                    return decodeURIComponent(cookiePair[1]);
                }
            }

            // Return null if not found
            return null;
        }
function init() {
            profileId= window.location.pathname.split("/").pop()


            console.log(profileId)
            $.ajax({
                type: "GET",
                contentType: "application/json",
                headers: {"Authorization":getCookie("Authorization")},
                url: "/auth/profile-info/" + profileId,
                error: function (response) {
                    alert("Get Profile Info Error " + JSON.parse(response.responseText).reason);
                    if (response.status === 405){
                        alert("You are not allowed to see this page:");
                        window.location.replace("/signin")
                    }
                }
            })

        }

        init()
    </script></head>
<body>
</body>
</html>
 
`