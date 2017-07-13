 //Google Plus Login
 var googleUser = {};
  var startApp = function() {
    gapi.load('auth2', function(){
      // Retrieve the singleton for the GoogleAuth library and set up the client.
      auth2 = gapi.auth2.init({
        client_id: '384146092934-7c8b1e18cg68617b6315a9eaq0dbd5o1.apps.googleusercontent.com',
        cookiepolicy: 'single_host_origin',
        // Request scopes in addition to 'profile' and 'email'
        //scope: 'additional_scope'
      });
      attachSignin(document.getElementById('singupGoogle'), document.getElementById('loginGoogle'));
    });
  };

  function attachSignin(elementSignup, elementLogin) {
    auth2.attachClickHandler(elementSignup, {},
        function(googleUser) {
			singupSocial(googleUser, "google", "signup");
        }, function(error) {
          console.log(JSON.stringify(error, undefined, 2));
        });
	auth2.attachClickHandler(elementLogin, {},
	function(googleUser) {
		singupSocial(googleUser, "google", "login");
	}, function(error) {
	  console.log(JSON.stringify(error, undefined, 2));
	});
  }
  
  
  //fb login
  var action;
  function statusChangeCallback(response) {
    console.log('statusChangeCallback');
    console.log(response);

    if (response.status === 'connected') {
      testAPI();
    } else {
      // The person is not logged into your app or we are unable to tell.

    }
  }



  window.fbAsyncInit = function() {
  FB.init({
    appId      : '291859141272771',
    cookie     : true,  // enable cookies to allow the server to access 
                        // the session
    xfbml      : true,  // parse social plugins on this page
    version    : 'v2.8' // use graph api version 2.8
  });

	$('#signupFB, #loginFB').on('click', function() {
			
		action = $(this).data('action');
		fbLogin();
	});
	
  };

function fbLogin() {
		FB.login(function(response) {
			if (response.status === 'connected') {
				testAPI();
			} else {
			// The person is not logged into this app or we are unable to tell. 
			}
		}, {scope: 'public_profile,email'});
	}

  (function(d, s, id) {
    var js, fjs = d.getElementsByTagName(s)[0];
    if (d.getElementById(id)) return;
    js = d.createElement(s); js.id = id;
    js.src = "//connect.facebook.net/en_US/sdk.js";
    fjs.parentNode.insertBefore(js, fjs);
  }(document, 'script', 'facebook-jssdk'));


  function testAPI() {
    console.log('Welcome!  Fetching your information.... ');
    FB.api('/me', function(response) {
		var fbUser = { token : FB.getAuthResponse()['accessToken'] }
		singupSocial(fbUser, "fb", action);
    });
  }
  
  startApp();
  
  
  
//submit social account data
function singupSocial(socialUser, type, action) {
	
	if (type ==="google") {
		user = socialUser.getBasicProfile();
		var user = {
			firstname: user.getGivenName(),
			lastname: user.getFamilyName(),
			email: user.getEmail(),
			image: user.getImageUrl(),
			kind: type,
			uid: user.getId(),
			token: socialUser.getAuthResponse().id_token,
			action: action
			
		}
	}
	if(type == "fb") {
		
		var user = {
			token: socialUser.token,
			kind: type,
			action: action
		}
	}
	console.log(user);
	$.post( "signupsocial",
		user,
		function( data ) {
			if (data.userAuth === true) {
				window.location.replace("/profile");
			}
		}, "json"
	)
	.fail(function (jqXHR, exception) {
			var msg = '';
			if (jqXHR.status === 0) {
				msg = 'Not connect.\n Verify Network.';
			} else if (jqXHR.status == 404) {
				msg = 'Requested page not found. [404]';
			} else if (jqXHR.status == 500) {
				msg = 'Internal Server Error [500].';
			} else if (exception === 'parsererror') {
				msg = 'Requested JSON parse failed.';
			} else if (exception === 'timeout') {
				msg = 'Time out error.';
			} else if (exception === 'abort') {0
				msg = 'Ajax request aborted.';
			} else {
				msg = 'Uncaught Error.\n' + jqXHR.responseText;
			}
			console.log(msg);
	});
}

//mnodal
$('#forgotPassword').on('show.bs.modal', function () {
  $('#login').modal('hide');
})

