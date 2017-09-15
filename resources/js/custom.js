var btnContinue;
var socket = new Ws("ws://bunity.com:8081/userchat");
	/*socket.OnConnect(function () {
		console.log("emit");
		socket.Emit("chat", "muie");
	});
	
	socket.On("chat", function (msg) {
		console.log("on");
	});
	
	socket.OnDisconnect(function () {
		console.log("disconnect");
	});*/
	socket.On("like", function (msg) {
		$(".notifications-count").removeClass("hidden").text(1);
		$.notify({
			// options
			message: 'Someone liked your page!' 
		},{
			// settings
			type: 'info',
			placement: {
				from: 'bottom',
				align: 'left'
			}
		});

	});
$(document).ready(function() {
	Dropzone.prototype.defaultOptions.dictRemoveFile = "";
	//$.fn.select2.defaults.set("theme", "classic");
	$('.businessCountry, .businessStateUsa, .businessStateCanada, .businessStateAustralia, .businessCateg, .businessCateg2').select2();
	
	$('.phonePrefix').select2({
		templateSelection: formatState,
		width: 'element',
	});
	loadBusinessGallery();
	loadBusinessesAjax();
	galleryUploadFunc();
	searchCartegoryList();
	$(".updateGallery").fancybox({
		afterClose : function() {
			$('#galleryUpload, #profileUpload, #coverUpload').remove();
		}
	});
	
	BusinessForm();
	BusinessFormSave();
	btnContinue = $(".businessContinue");
	btnContinue.on("click", function(e) {
		var href = this.href;
		e.preventDefault();
		ajaxBusiness(href, $(this));
	});
	
	$(".businessSendNumber").on("click", function(e) {
		var hrefSMS = this.href;
		e.preventDefault();
		$.post( "/businesses/sendsms", {"prefix": $(".phonePrefix").val(), "smsCode":$(".phoneSMSField").val(), "businessID": $("#businessID").val()},
			function( data ) {
				if(data.constructor === Array && data.length > 0) {
						$(".phoneSMSField").addClass('error-form').attr('data-original-title', data[0].Message)
							.tooltip({trigger:'hover'}).tooltip('fixTitle').tooltip('show');
						$( ".phoneSMSField" ).focus(function() {
							$( this ).removeClass("error-form").attr('data-original-title', '').tooltip('hide');
						});
				} else {
					window.location = hrefSMS;
				}
		}, "json")
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
	});
	
	$(".businessVerifyCode").on("click", function(e) {
		e.preventDefault();
		$.post( "/businesses/verifycode", {"verificationCode":$(".verificationCode").val(), "businessID": $("#businessID").val()},
			function( data ) {
				if(!data.response) {
					$(".verificationCode").addClass('error-form').attr('data-original-title', data.message)
						.tooltip({trigger:'hover'}).tooltip('fixTitle').tooltip('show');
					$( ".verificationCode" ).focus(function() {
						$( this ).removeClass("error-form").attr('data-original-title', '').tooltip('hide');
					});
				}
				else {
					window.location = "/businesses";
				}
		}, "json")
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
	});
	
	$(".updateGallery").on("click", function(e) {
		e.preventDefault();
		var actionType = $(this).data("action");
		$.post( "/businesses/updatephotos", {"userID":$("#userID").val(), "businessID": $("#businessID").val()},
			function( data ) {
				thumbnailUrls=[];
				profileThumbnailUrls=[];
				coverThumbnailUrls=[];
				
				if(data.galleryImages) {
					for (var i = 0; i < data.galleryImages.length; i++) {
						thumbnailUrls.push(data.galleryImages[i]);
					}
				}
				
				if(data.profileImages) {
					for (var i = 0; i < data.profileImages.length; i++) {
						profileThumbnailUrls.push(data.profileImages[i]);
					}
				}
				
				if(data.coverImages) {
					for (var i = 0; i < data.coverImages.length; i++) {
						coverThumbnailUrls.push(data.coverImages[i]);
					}
				}
				if(actionType == "gallery") {
					$('#imageUploadContainer').html('<button data-fancybox-close="" class="fancybox-close-small" title="Close"></button><div class="editPhotoTitle">Edit Gallery Photos</div><div class="editPhotoSubTitle">Choose up to 8 photos you\'d like to feature.</div><div id="galleryUpload" class="dropzone"><div class="dz-default dz-message"><div class="dropzoneBox"><i class="fa fa-camera" aria-hidden="true"></i></div><div class="dropzoneBox"><i class="fa fa-camera" aria-hidden="true"></i></div><div class="dropzoneBox"><i class="fa fa-camera" aria-hidden="true"></i></div><div class="dropzoneBox"><i class="fa fa-camera" aria-hidden="true"></i></div><div class="dropzoneBox" style="clear:both;"><i class="fa fa-camera" aria-hidden="true"></i></div><div class="dropzoneBox"><i class="fa fa-camera" aria-hidden="true"></i></div><div class="dropzoneBox"><i class="fa fa-camera" aria-hidden="true"></i></div><div class="dropzoneBox"><i class="fa fa-camera" aria-hidden="true"></i></div></div></div>');
				} else if(actionType == "profile") {
					$('#imageUploadContainer').html('<button data-fancybox-close="" class="fancybox-close-small" title="Close"></button><div class="editPhotoTitle">Edit Profile Photo</div><div class="editPhotoSubTitle">Choose 1 photo you\'d like to feature.</div><div id="profileUpload" class="dropzone"><div class="dz-default dz-message"><div class="dropzoneBox"><i class="fa fa-camera" aria-hidden="true"></i></div></div></div>');
				} else if(actionType == "cover") {
					$('#imageUploadContainer').html('<button data-fancybox-close="" class="fancybox-close-small" title="Close"></button><div class="editPhotoTitle">Edit Cover Photo</div><div class="editPhotoSubTitle">Choose 1 photo you\'d like to feature.</div><div id="coverUpload" class="dropzone"><div class="dz-default dz-message"><div class="dropzoneBox"><i class="fa fa-camera" aria-hidden="true"></i></div></div></div>');
				}
				galleryUploadFunc();
		}, "json")
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
	});
	
	var likeFlag = false;
	$(".businessLikeBtn").on("click", function(e) {
		e.preventDefault();
		if(!likeFlag) {
			var thisbtn = $(this);
			likeFlag = true;
			var nrLikes = parseInt($(".nrLikes").text());
			var likesBefore = nrLikes;
			if(liked) {
				--nrLikes;
				thisbtn.removeClass("liked");
			} else {
				++nrLikes;
				thisbtn.addClass("liked");
			}
			$(".nrLikes").html(nrLikes);
			
			e.preventDefault();
			var actionType = $(this).data("action");
			var businessID = $("#businessID").val();
			$.post( "/likes/"+businessID,
				function( data ) {
					if(data.success === true) {
						liked = true;
						thisbtn.addClass("liked");
						socket.OnConnect(function () {
							socket.Emit("like", userID);
						});
					} else {
						liked = false;
						thisbtn.removeClass("liked");
					}
					$(".nrLikes").html(data.count);
					likeFlag = false;
			}, "json")
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
					$(".nrLikes").html(likesBefore);
					if(liked) {
						thisbtn.addClass("liked");
					} else {
						thisbtn.removeClass("liked");
					}
					likeFlag = false;
					console.log(msg);
			});
		}
	});
	
	if($("#businessDescription").length > 0) {
		tinymce.init({
			selector: '#businessDescription',
			branding: false,
			height: 500,
			paste_as_text: true,
			menubar: false,
			plugins: [
				'paste lists',
				'wordcount'
			],
			toolbar: 'undo redo | insert | styleselect | bold italic | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent | link image',
			setup:function(ed) {
				var timeoutId;
				ed.on('keyup', function(e) {
					clearTimeout(timeoutId);
					timeoutId = setTimeout(function() {
						//console.log(ed.getContent());
						
						ajaxBusiness();
					}, 750);
				   
				});
				ed.on('change', function(e) {
						ajaxBusiness();
				});

			}
		});
	}
	
	$(document).click(function (e) {
		var container = $(".headSearch");
		if (!container.is(e.target))
		{
			$('.searchResults').hide();
		}
		if(!$(".filterCateg").is(e.target) && !$(".choose_categoryBtn").is(e.target)) {
			$(".choose_categoryBtn").show();
			$(".filterCateg").addClass("hidden");
			$(".filterCateg").val('');
			$(".categFilterList li").removeAttr("style");
		}
	});
	
	var currentRequest = null; 
	var headSearch = $(".headSearch");
	//var timeoutId2;
	
	
	headSearch.on("keyup focus",function() {
		
		var str = headSearch.val();
		$(".searchResults").html('<ul class="list-unstyled"><li>'+str.toLowerCase()+'</li></ul><div class="searchLoad">FINDING RESULTS <img color="white" src="/static/images/LOOn0JtHNzb.gif" class="img" alt="" width="16" height="16"></div>');
		var strArr = str.split(" ");
		var tempStr = 0;
		for(var i=0; i< strArr.length; i++) {
			if(strArr[i].length > tempStr) {
				tempStr = strArr[i].length;
			}
		}
		if(tempStr >= 3) {
		
		
			//clearTimeout(timeoutId2);
			//timeoutId2 = setTimeout(function() {
			
			currentRequest = $.ajax({
			type:"POST",
			url:"/livesearch",
			data: {"keyword":str},
			beforeSend : function()    {       
				if(currentRequest != null) {
					currentRequest.abort();
					console.log(1);
				}
			},
			success: function(data) {
				console.log("succeess");
			    if(data.results.length > 0) {
					var html = '<ul class="list-unstyled">';
					for(var i=0;i<data.results.length;i++) {
						var strL = str.toLowerCase();
						var n = data.results[i].Name.replace(strL/*new RegExp(strL, 'g')*/, "<strong>"+strL+"</strong>");
						var b = data.results[i];
						html += '<li class=""><a href="/'+data.results[i].Url+'"><div class="pull-left imgLiveSearch"><img src="/static/uploads/'+b.UserId+'/'+b.Url+'/profile/'+b.Image[0]+'"></div><div class="liveSearchRIght"><div>'+n+'</div><div><small class="liveSearchCat">'+b.Category+'</small></div></div></a></li>';
					}
					html += '</ul>';
				} else {
					$(".searchLoad").remove();
				}
				$(".searchResults").html(html).show();
			},
			dataType: 'json',
		  });
		//}, 400);
		/*$.post( "/livesearch", {"keyword":str},
		function( data ) {
			if(data.results.length > 0) {
				var html = '<ul class="list-unstyled">';
				for(var i=0;i<data.results.length;i++) {
					var strL = str.toLowerCase();
					var n = data.results[i].Name.replace(strL/*new RegExp(strL, 'g')*//*, "<strong>"+strL+"</strong>");
					var b = data.results[i];
					html += '<li class=""><a href="/'+data.results[i].Url+'"><div class="pull-left imgLiveSearch"><img src="/static/uploads/'+b.UserId+'/'+b.Url+'/profile/'+b.Image[0]+'"></div><div class="liveSearchRIght"><div>'+n+'</div><div><small class="liveSearchCat">'+b.Industry+'</small></div></div></a></li>';
				}
				html += '</ul>';
			} else {
				$(".searchLoad").remove();
			}
			$(".searchResults").html(html).show();
		}, "json")
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
		});*/
		} else {
			$(".searchLoad").remove();
		}
	});
	
});



var BusinessFormSave = function() {
	var timeoutId;
	$('.businessForm input, .businessForm select').not('.phonePrefixC select').on('input select2:select', function() {
		var _this = $(this);
		var nodeName = $(this).prop('nodeName');
		if(nodeName === "INPUT") {
			clearTimeout(timeoutId);
			timeoutId = setTimeout(function() {
				if(_this.val() != "" || _this.hasClass("businessSocial") || _this.hasClass("businessWebsite") || _this.hasClass("businessSocial") || $(".businessState").is(':enabled')) {
					ajaxBusiness();
				}
			}, 750);
		} else if(nodeName === "SELECT") {
			if(_this.val() != "" || _this.hasClass("businessSocial") || _this.hasClass("businessWebsite")) {
				ajaxBusiness();
			}
		}
	});
}

var ajaxBusiness = function(url, btn) {
	var formSelector = $(".businessForm");
	var values = {};
		values.add = 0;
		values.back = 0;
	if($("#businessID").val()) {
		values.businessID = $("#businessID").val();
	}
	if (btn !== undefined) {
		values.add = btn.data("add");
		values.back = btn.data("back");
	}
	
	$.each($(formSelector).serializeArray(), function(i, field) {
			values[field.name] = field.value;
	});
	if($(".formStep").val() === "3") {
		values["business[description]"] = tinymce.activeEditor.getContent();
	}
	if($(".businessID").val()) {
		values["businessID]"] = $(".businessID").val();
	}
	$.post( "/businesses/trackEvents", values,
			function( data ) {
				if (url !== undefined) {
					if(data.constructor === Array && data.length > 0) {
						
						for(i=0; i < data.length; i++) {
							$("."+data[i].Class).not(':hidden').addClass('error-form').attr('data-original-title', data[i].Message)
								.tooltip({trigger:'hover'}).tooltip('fixTitle').tooltip('show');
						}
						$( ".businessForm :input" ).focus(function() {
							$( this ).removeClass("error-form").attr('data-original-title', '').tooltip('hide');
						});

					}
					else {		
						if(values.add == 1) {
							window.location = "/businesses/add/step3/"+data;
							return;
						}
						window.location = url;
					}
				}
				var d = new Date();
				$('.form-status-holder').html('<strong>Saved! Last: ' + d.toLocaleTimeString()+"</strong>");
		}, "json")
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

var BusinessForm = function() {
	var stateSelect = $(".businessState");
	var stateUsaSelect = $(".businessStateUsa");
	var stateCanadaSelect = $(".businessStateCanada");
	var stateAustraliaSelect = $(".businessStateAustralia");
	
	var select2Australia = $(".statesAustraliaSelect");
	var select2Usa = $(".statesUsaSelect");
	var select2Canada = $(".statesCanadaSelect");
	
	var countrySelect = $(".businessCountry");
	var flag = false;
	var currentCountry = countrySelect.val();
	var oldCountry = currentCountry;
	if (oldCountry == "Australia" || oldCountry == "United States" || oldCountry == "Canada") {
		flag = true;
	}
	countrySelect.on("select2:select", function(){
		oldCountry = currentCountry;
		$( ".businessStateControl" ).removeClass("error-form").attr('data-original-title', '').tooltip('hide');
		currentCountry = $(this).val();
		if(currentCountry == 'United States'){
			flag = true;
			stateSelect.prop('disabled', true).addClass("hidden");
			stateCanadaSelect.prop('disabled', true).addClass("hidden");
			select2Canada.addClass("hidden");
			stateAustraliaSelect.prop('disabled', true).addClass("hidden");
			select2Australia.addClass("hidden");
			stateUsaSelect.prop('disabled', false).removeClass("hidden");
			select2Usa.removeClass("hidden");
		} else if(currentCountry == 'Canada') {
			flag = true;
			stateSelect.prop('disabled', true).addClass("hidden");
			stateUsaSelect.prop('disabled', true).addClass("hidden");
			select2Usa.addClass("hidden");
			stateAustraliaSelect.prop('disabled', true).addClass("hidden");
			select2Australia.addClass("hidden");
			stateCanadaSelect.prop('disabled', false).removeClass("hidden");
			select2Canada.removeClass("hidden");
		} else if(currentCountry == 'Australia') {
			flag = true;
			stateSelect.prop('disabled', true).addClass("hidden");
			stateUsaSelect.prop('disabled', true).addClass("hidden");
			select2Usa.addClass("hidden");
			stateCanadaSelect.prop('disabled', true).addClass("hidden");
			select2Canada.addClass("hidden");
			stateAustraliaSelect.prop('disabled', false).removeClass("hidden");
			select2Australia.removeClass("hidden");
		} else {
			flag = false;
			stateCanadaSelect.prop('disabled', true).addClass("hidden");
			select2Canada.addClass("hidden");
			stateUsaSelect.prop('disabled', true).addClass("hidden");
			select2Usa.addClass("hidden");
			stateAustraliaSelect.prop('disabled', true).addClass("hidden");
			select2Australia.addClass("hidden");
			stateSelect.prop('disabled', false).removeClass("hidden");
		}
		console.log(oldCountry);
		console.log(flag);
		if (flag && (oldCountry != "Australia" && oldCountry != "United States" && oldCountry != "Canada")) {
			$(".businessStateControl ").val("");
		}
		if (!flag && (oldCountry == "Australia" || oldCountry == "United States" || oldCountry == "Canada")) {
			$(".businessStateControl ").val("");
		}
	});	
}

function galleryUploadFunc() {
	Dropzone.autoDiscover = false;
	var fileList = new Array;
	var fileList2 = new Array;
	var fileList3 = new Array;
	$("#galleryUpload").dropzone({
		url: "/businesses/addfiles",
		sending: function(file, xhr, formData){
			formData.append('imageType', "gallery");
            formData.append('businessID', $("#businessID").val());
			formData.append('imageFormat', file.type);
			console.log(file.type);
        },
		addRemoveLinks : true,
		maxFiles:8,
		acceptedFiles: ".jpeg,.jpg,.png",
		init: function() {
			var myDropzone = this;
			var existingFileCount = thumbnailUrls.length;
			//myDropzone.options.maxFiles = myDropzone.options.maxFiles - existingFileCount ;
			if (thumbnailUrls) {
				for (var i = 0; i < thumbnailUrls.length; i++) {
					var imgURL = "/static/uploads/"+$("#userID").val()+"/"+$("#businessID").val()+"/gallery/"+thumbnailUrls[i];
					var mockFile = { 
						name: thumbnailUrls[i], 
						//size: 12345, 
						//type: 'image/jpeg', 
						status: Dropzone.ADDED, 
						url: imgURL,
						accepted: true
					};

					// Call the default addedfile event handler
					myDropzone.emit("addedfile", mockFile);
					myDropzone.emit("complete", mockFile);
					// And optionally show the thumbnail of the file:
					myDropzone.emit("thumbnail", mockFile, imgURL);

					myDropzone.files.push(mockFile);
				}
				console.log(myDropzone.options.maxFiles);
				console.log(myDropzone.files.length);
			}
		},
	error: function(file, message, xhr) {
		 //$(file.previewElement).remove();
		newServerName = "";
		this.removeFile(file);
	},
	success: function(file, serverFileName) {
		
		fileList.push ({"serverFileName" : serverFileName.fname, "fileName" : file.name});
		var lastChar = serverFileName.fname.charAt(serverFileName.fname.length - 5);
		if(lastChar == 1) {
			var orientation = "landscape";
		} else {
			var orientation = "portrait";
		}
		$(".galleryContainer").append('<div class="galleryImageSIngle '+orientation+'"><a class="businessGallery" data-fancybox="gallery" href="'+serverFileName.url+'"><img src="'+serverFileName.url+'"></a></div>');
		
		loadBusinessGallery();
	},
	removedfile: function(file) {
		var rmvFile = "";
		console.log("delete file length: "+fileList.length);
		if(fileList.length > 0) {
		for(f=0;f<fileList.length;f++){

			if(fileList[f].fileName == file.name)
			{
				console.log("new: "+ fileList[f].fileName+"old: "+file.name);
				rmvFile = fileList[f].serverFileName;
				fileList.splice(f,1);
				//myDropzone.options.maxFiles = myDropzone.options.maxFiles + 1;
			}

		}
		} else {
			rmvFile = file.name;
		}
		if(rmvFile == "") {
			rmvFile = file.name;
		}
		if (rmvFile){
			//console.log(rmvFile);
			$.ajax({
				type: 'POST',
				url: '/businesses/deletefile',
				data: "id="+rmvFile+"&businessID="+$("#businessID").val()+"&imageType=gallery",
				dataType: 'html'
			}).done(function () { $(document).find(file.previewElement).remove(); $('.galleryContainer img[src*="'+rmvFile+'"]').closest(".galleryImageSIngle").remove(); loadBusinessGallery(); });
			//var _ref;
			
			//return (_ref = file.previewElement) != null ? _ref.parentNode.removeChild(file.previewElement) : void 0;   
		}
		//console.log(this.options.maxFiles);
	}
	});
	
	var $cropperModal = $(modalTemplate);
	$("#profileUpload").dropzone({
		url: "/businesses/addfiles",
		sending: function(file, xhr, formData){
			formData.append('imageType', "profile");
            formData.append('businessID', $("#businessID").val());
			formData.append('imageFormat', file.type);
			//console.log(file.type);
        },
		addRemoveLinks : true,
		maxFiles:2,
		maxFilesize: 5,
		acceptedFiles: ".jpeg,.jpg,.png",
		autoProcessQueue : false,
		autoQueue: false,
		accept: function(file, done) {
			if(file.cropped) {
				done();
				console.log("hai");
			}
			file.acceptDimensions = done;
            file.rejectDimensions = function(msg) { done(msg); };
		},
		init: function() {
		
			var myDropzone = this;
			if (profileThumbnailUrls) {
				for (var i = 0; i < profileThumbnailUrls.length; i++) {
					var imgURL = "/static/uploads/"+$("#userID").val()+"/"+$("#businessID").val()+"/profile/"+profileThumbnailUrls[i];
					var mockFile = { 
						name: profileThumbnailUrls[i], 
						//size: 12345, 
						//type: 'image/jpeg', 
						status: Dropzone.ADDED, 
						url: imgURL
					};

					// Call the default addedfile event handler
					myDropzone.emit("addedfile", mockFile);
					myDropzone.emit("complete", mockFile);
					// And optionally show the thumbnail of the file:
					myDropzone.emit("thumbnail", mockFile, imgURL);

					myDropzone.files.push(mockFile);
					var existingFileCoun = profileThumbnailUrls.length;
					myDropzone.options.maxFiles = myDropzone.options.maxFiles - existingFileCoun;
				}
			}
		},
	success: function(file, serverFileName) {
		
		fileList2.push ({"serverFileName" : serverFileName.fname, "fileName" : file.name});
		var lastChar = serverFileName.fname.charAt(serverFileName.fname.length - 5);
		if(lastChar == 1) {
			var orientation = "landscape";
		} else {
			var orientation = "portrait";
		}
		$(".profile-picture").removeClass("landscape portrait").addClass(orientation);
		$(".profile-picture img").attr("src", serverFileName.url);
		$.fancybox.close();
		if(serverFileName.fname) {
			$cropperModal.modal('hide').remove();
			$(".modal-backdrop").remove();
		}
	},
	removedfile: function(file) {
			var rmvFile = "";
		//console.log("delete file length: "+fileList2.length);
		if(fileList2.length > 0) {
		for(f=0;f<fileList2.length;f++){

			if(fileList2[f].fileName == file.name)
			{
				//console.log("new: "+ fileList2[f].fileName+"old: "+file.name);
				rmvFile = fileList2[f].serverFileName;
				fileList2.splice(f,1);
				//myDropzone.options.maxFiles = myDropzone.options.maxFiles + 1;
			}

		}
		} else {
			rmvFile = file.name;
		}
		if(rmvFile == "") {
			rmvFile = file.name;
		}
		if (rmvFile){
			this.options.maxFiles = 1;
			//console.log(rmvFile);
			$.ajax({
				type: 'POST',
				url: '/businesses/deletefile',
				data: "id="+rmvFile+"&businessID="+$("#businessID").val()+"&imageType=profile",
				dataType: 'html'
			}).done(function () { $(document).find(file.previewElement).remove(); });
			//var _ref;
			
			//return (_ref = file.previewElement) != null ? _ref.parentNode.removeChild(file.previewElement) : void 0;   
		}
		//console.log(this.options.maxFiles);
		},
	thumbnail: function(file) {
		if (file.acceptDimensions) {
			if (file.width < 160 && file.height < 160) {
				file.rejectDimensions("The resolution has to be at least 160 x 160");
			}	else {

				file.acceptDimensions();
			}
		}
		$(".dz-image img").attr("src",file.url)
		if(file.accepted) {
			$.fancybox.close();
			var myDropzone = this
			if (file.cropped) {
				return;
			}
			var cachedFilename = file.name;
			//console.log(file);
			//myDropzone.removeFile(file);
		
			
			var $uploadCrop = $cropperModal.find('.crop-upload');
			var $img = $('<img />');
			var reader = new FileReader();
			reader.onloadend = function () {
				$cropperModal.find('.image-container').html($img);
				$img.attr('src', reader.result);
				$img.cropper({
					preview: '.image-preview',
					aspectRatio: 1 / 1,
					autoCropArea: 1,
					movable: false,
					cropBoxResizable: true,
					minContainerHeight : 320,
					minContainerWidth : 568,
					viewMode:2,
					minCropBoxHeight: 160,
					minCropBoxWidth:160
				});
			};
			
			reader.readAsDataURL(file);		
			$cropperModal.modal('show');
				
			$uploadCrop.on('click', function() {
				var blob = $img.cropper('getCroppedCanvas').toDataURL();
				var newFile = dataURItoBlob(blob);
				newFile.cropped = true;
				newFile.name = cachedFilename;

				myDropzone.removeAllFiles();
				myDropzone.options.maxFiles = 1;
				
				myDropzone.addFile(newFile);
				myDropzone.enqueueFile(newFile);
				myDropzone.processQueue();
				
				
			});
			 var $this = $(document);
			$this.on('click', '.rotate-right', function () {
                $img.cropper('rotate', 45);
            })
            .on('click', '.rotate-left', function () {
                $img.cropper('rotate', -45);
            })
            .on('click', '.reset', function () {
                $img.cropper('reset');
            })
            .on('click', '.scale-x', function () {
                var $this = $(this);
                $img.cropper('scaleX', $this.data('value'));
                $this.data('value', -$this.data('value'));
            })
            .on('click', '.scale-y', function () {
                var $this = $(this);
                $img.cropper('scaleY', $this.data('value'));
                $this.data('value', -$this.data('value'));
            })
			
			.on('click', '.zoom-in', function () {
				$img.cropper('zoom', -0.1);
            })
			
			.on('click', '.zoom-out', function () {
                $img.cropper('zoom', 0.1);
            });
		}
	}
	});
	

	$("#coverUpload").dropzone({
		url: "/businesses/addfiles",
		sending: function(file, xhr, formData){
			formData.append('imageType', "cover");
            formData.append('businessID', $("#businessID").val());
			formData.append('imageFormat', file.type);
			console.log(file.type);
        },
		addRemoveLinks : true,
		maxFiles:2,
		autoProcessQueue : false,
		autoQueue: false,
		accept: function(file, done) {
			if(file.cropped) {
				done();
			}
			file.acceptDimensions = done;
            file.rejectDimensions = function(msg) { done(msg); };
		},
		acceptedFiles: ".jpeg,.jpg,.png",
		init: function() {
			var myDropzone = this;
			if (coverThumbnailUrls) {
				for (var i = 0; i < coverThumbnailUrls.length; i++) {
					var imgURL = "/static/uploads/"+$("#userID").val()+"/"+$("#businessID").val()+"/cover/"+coverThumbnailUrls[i];
					var mockFile = { 
						name: coverThumbnailUrls[i], 
						//size: 12345, 
						//type: 'image/jpeg', 
						status: Dropzone.ADDED, 
						url: imgURL
					};

					// Call the default addedfile event handler
					myDropzone.emit("addedfile", mockFile);
					myDropzone.emit("complete", mockFile);
					// And optionally show the thumbnail of the file:
					myDropzone.emit("thumbnail", mockFile, imgURL);

					myDropzone.files.push(mockFile);
					var existingFileCoun = coverThumbnailUrls.length;
					myDropzone.options.maxFiles = myDropzone.options.maxFiles - existingFileCoun;
				}
			}
		}
		,
	success: function(file, serverFileName) {
		
		fileList3.push ({"serverFileName" : serverFileName.fname, "fileName" : file.name});
		$(".head-banner img").attr("src", serverFileName.url);
		$.fancybox.close();
		if(serverFileName.fname) {
			$cropperModal.modal('hide').remove();
			$(".modal-backdrop").remove();
		}
	},
	removedfile: function(file) {
		var rmvFile = "";
		if(fileList3.length > 0) {
		for(f=0;f<fileList3.length;f++){

			if(fileList3[f].fileName == file.name)
			{
				console.log("new: "+ fileList3[f].fileName+"old: "+file.name);
				rmvFile = fileList3[f].serverFileName;
				fileList3.splice(f,1);
				//myDropzone.options.maxFiles = myDropzone.options.maxFiles + 1;
			}

		}
		} else {
			rmvFile = file.name;
		}
		if(rmvFile == "") {
			rmvFile = file.name;
		}
		if (rmvFile){
			this.options.maxFiles = 1;
			console.log(rmvFile);
			$.ajax({
				type: 'POST',
				url: '/businesses/deletefile',
				data: "id="+rmvFile+"&businessID="+$("#businessID").val()+"&imageType=cover",
				dataType: 'html'
			}).done(function () { $(document).find(file.previewElement).remove(); });
			//var _ref;
			
			//return (_ref = file.previewElement) != null ? _ref.parentNode.removeChild(file.previewElement) : void 0;   
		}
		console.log(this.options.maxFiles);
	},
	thumbnail: function(file) {
		if (file.acceptDimensions) {
			if (file.width < 840 && file.height < 285) {
				file.rejectDimensions("The resolution has to be at least 840 x 285");
			}	else {
				file.acceptDimensions();
			}
		}
		$(".dz-image img").attr("src",file.url)
		if(file.accepted) {
			$.fancybox.close();
			var myDropzone = this
			if (file.cropped) {
				return;
			}
			var cachedFilename = file.name;
			//console.log(file);
			//myDropzone.removeFile(file);
		
			
			var $uploadCrop = $cropperModal.find('.crop-upload');
			$cropperModal.find(".zoom-in, .zoom-out").remove();
			var $img = $('<img />');
			var reader = new FileReader();
			reader.onloadend = function () {
				$cropperModal.find('.image-container').html($img);
				$img.attr('src', reader.result);
				$img.cropper({
					preview: '.image-preview',
					aspectRatio: 35  / 12,
					autoCropArea: 1,
					movable: false,
					cropBoxResizable: true,
					minContainerHeight : 320,
					minContainerWidth : 568,
					viewMode:2,
					cropBoxResizable: false
				});
			};
			
			reader.readAsDataURL(file);		
			$cropperModal.modal('show');
				
			$uploadCrop.on('click', function() {
				var blob = $img.cropper('getCroppedCanvas').toDataURL();
				var newFile = dataURItoBlob(blob);
				newFile.cropped = true;
				newFile.name = cachedFilename;

				  myDropzone.removeAllFiles();
				  myDropzone.options.maxFiles = 1;
		
				myDropzone.addFile(newFile);
				myDropzone.enqueueFile(newFile);
				myDropzone.processQueue();
				
				
			});
			 var $this = $(document);
			$this.on('click', '.rotate-right', function () {
                $img.cropper('rotate', 45);
            })
            .on('click', '.rotate-left', function () {
                $img.cropper('rotate', -45);
            })
            .on('click', '.reset', function () {
                $img.cropper('reset');
            })
            .on('click', '.scale-x', function () {
                var $this = $(this);
                $img.cropper('scaleX', $this.data('value'));
                $this.data('value', -$this.data('value'));
            })
            .on('click', '.scale-y', function () {
                var $this = $(this);
                $img.cropper('scaleY', $this.data('value'));
                $this.data('value', -$this.data('value'));
            });
			
			/*.on('click', '.zoom-in', function () {
				$img.cropper('zoom', -0.1);
            })
			
			.on('click', '.zoom-out', function () {
                $img.cropper('zoom', 0.1);
            });*/
		}
	}
	});
}

function loadBusinessGallery() {
	$(".businessGallery").fancybox({
		openEffect	: 'fade',
		closeEffect	: 'fade',
		type : "image",
		thumbs : {
			autoStart : true,
			hideOnClose : true
		}
	});
}

function insensitiveReplaceAll(original, find, replace) {
  var str = "",
    remainder = original,
    lowFind = find.toLowerCase(),
    idx;

  while ((idx = remainder.toLowerCase().indexOf(lowFind)) !== -1) {
    str += remainder.substr(0, idx) + replace;

    remainder = remainder.substr(idx + find.length);
  }

  return str + remainder;
}

function stringPos(textStr) {
	return 
	extStr.charAt(textStr.length-5);
}

function loadBusinessesAjax() {
	var countPage = 2;
	if($(".searchAjax").length) {
		var win = $(window);
		// Each time the user scrolls
		console.log($(document).height() - win.height());
			console.log(win.scrollTop());
		win.scroll(function() {
			// End of the document reached?
			console.log($(document).height() - win.height());
			console.log(win.scrollTop());
			if ($(document).height() - win.height() == win.scrollTop()) {
				$('#loading').show();

				$.post( "/search/business", {"q": $(".searchTerm").val(), "countPage": countPage, "business_category": $(".businessCategory").val(), "verified": $(".verifiedField").val()},
				function( data ) {
					var b = data.businesses;
					if(b) {
						var html = '';
						for (var i=0; i < b.length; i++) {
							html += `
							<li>
								<a class="bProfilePic" href="/`+b[i]._id.toString()+`" class="">
									<img width="72" height="72" src="/static/uploads/`+b[i].user_id.toString()+`/`+b[i]._id.toString()+`/profile/`+b[i].profile[0]+`">
								</a>
								<div class="">
									<a href="/`+b[i]._id.toString()+`" class="">`+b[i].name+`</a>
									<p style="margin-bottom:0;">`+b[i].categ+`</p>
									<p style="margin-bottom:0;">`+b[i].address.city+`, `+b[i].address.country+`</p>
									<p style="margin-bottom:0;">`+b[i].nrLikes+` like this</p>
								</div>
							</li>`;
						}
						$(".bizFindResults").append(html);
						countPage++;
					}
				
				}, "json")
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
		});
	}
}

function searchCartegoryList() {
	$(".choose_categoryBtn").on("click", function() {
		$(this).hide();
		$(".filterCateg").removeClass("hidden");
		$(".filterCateg")[0].focus();
	});
	
	var liFilter = $('.categFilterList li');
	$('.filterCateg').on('keyup', function () {
		var value = this.value;
		
	   liFilter.hide().each(function () {
			if ($(this).find("a").text().toLowerCase().search(value.toLowerCase()) > -1) {
				$(this).attr("style", "display: list-item !important");
			}
			if (value == "" && !$(this).find("a").hasClass("active")) {
				 $(this).removeAttr("style");
			}
		});
	});
}

function formatState (state) {
    return state.id;
	
};


var modalTemplate = '' + 
	'<div class="modal fade" tabindex="-1" role="dialog">' + 
		'<div class="modal-dialog" role="document">' + 
			'<div class="modal-content">' + 
				'<div class="modal-header">' + 
					'<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>' + 
					'<h4 class="modal-title">Crop Image</h4>' + 
				'</div>' + 						
				'<div class="modal-body">' + 
					'<div class="image-container"></div>' + 
				'</div>' + 						
				'<div class="modal-footer">' + 
					'<button type="button" class="btn btn-warning zoom-in"><span class="fa fa-search-minus"></span></button>' +
					'<button type="button" class="btn btn-warning zoom-out"><span class="fa fa-search-plus"></span></button>' +
					
					'<button type="button" class="btn btn-warning rotate-left"><span class="fa fa-rotate-left"></span></button>' +
					'<button type="button" class="btn btn-warning rotate-right"><span class="fa fa-rotate-right"></span></button>' +
					'<button type="button" class="btn btn-warning scale-x" data-value="-1"><span class="fa fa-arrows-h"></span></button>' +
					'<button type="button" class="btn btn-warning scale-y" data-value="-1"><span class="fa fa-arrows-v"></span></button>' +
					'<button type="button" class="btn btn-warning reset"><span class="fa fa-refresh"></span></button>' +
					'<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>' + 
					'<button type="button" class="btn btn-primary crop-upload">Upload</button>' + 
				'</div>' + 
			'</div>' + 
		'</div>' + 
	'</div>' + 
'';

function dataURItoBlob(dataURI) {
    var byteString;
    if (dataURI.split(',')[0].indexOf('base64') >= 0)
        byteString = atob(dataURI.split(',')[1]);
    else
        byteString = unescape(dataURI.split(',')[1]);

    // separate out the mime component
    var mimeString = dataURI.split(',')[0].split(':')[1].split(';')[0];

    // write the bytes of the string to a typed array
    var ia = new Uint8Array(byteString.length);
    for (var i = 0; i < byteString.length; i++) {
        ia[i] = byteString.charCodeAt(i);
    }

    return new Blob([ia], {type:mimeString});
}



/*var messageTxt;
var messages;
$(function () {
	messageTxt = $("#messageTxt");
	messages = $("#messages");
	w = new WebSocket("ws://" + HOST + "/userchat");
	w.onopen = function () {
		console.log("Websocket connection enstablished");
	};
	w.onclose = function () {
		appendMessage($("<div><center><h3>Disconnected</h3></center></div>"));
	};
	w.onmessage = function(message){
		console.log(message.data);
		appendMessage($("<div>" + message.data + "</div>"));
	};
	$("#sendBtn").click(function () {
		w.send(messageTxt.val().toString());
		messageTxt.val("");
	});
})
function appendMessage(messageDiv) {
    var theDiv = messages[0];
    var doScroll = theDiv.scrollTop == theDiv.scrollHeight - theDiv.clientHeight;
    messageDiv.appendTo(messages);
    if (doScroll) {
        theDiv.scrollTop = theDiv.scrollHeight - theDiv.clientHeight;
    }
}*/