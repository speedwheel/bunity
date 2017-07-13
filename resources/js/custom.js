var btnContinue;
$(document).ready(function() {
	loadBusinessGallery();
	
	galleryUploadFunc();
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
		$.post( "/businesses/sendsms", {"smsCode":$(".phoneSMSField").val(), "businessID": $("#businessID").val()},
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
				console.log(thumbnailUrls);
				
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
					$('#imageUploadContainer').append('<div id="galleryUpload" class="dropzone"></div>');
				} else if(actionType == "profile") {
					$('#imageUploadContainer').append('<div id="profileUpload" class="dropzone"></div>');
				} else if(actionType == "cover") {
					$('#imageUploadContainer').append('<div id="coverUpload" class="dropzone"></div>');
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
});



var BusinessFormSave = function() {
	var timeoutId;
	$('.businessForm input, .businessForm select').on('input', function() {
		var _this = $(this);
		var nodeName = $(this).prop('nodeName');
		if(nodeName === "INPUT") {
			clearTimeout(timeoutId);
			timeoutId = setTimeout(function() {
				if(_this.val() != "" || _this.hasClass("businessSocial") || _this.hasClass("businessWebsite")) {
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
	var countrySelect = $(".businessCountry");
	var flag = false;
	var currentCountry = countrySelect.val();
	var oldCountry = currentCountry;
	if (oldCountry == "Australia" || oldCountry == "United States" || oldCountry == "Canada") {
		flag = true;
	}
	countrySelect.change(function(){
		oldCountry = currentCountry;
		$( ".businessStateControl" ).removeClass("error-form").attr('data-original-title', '').tooltip('hide');
		currentCountry = $(this).val();
		if(currentCountry == 'United States'){
			flag = true;
			stateSelect.prop('disabled', true).addClass("hidden");
			stateCanadaSelect.prop('disabled', true).addClass("hidden");
			stateAustraliaSelect.prop('disabled', true).addClass("hidden");
			stateUsaSelect.prop('disabled', false).removeClass("hidden");
		} else if(currentCountry == 'Canada') {
			flag = true;
			stateSelect.prop('disabled', true).addClass("hidden");
			stateUsaSelect.prop('disabled', true).addClass("hidden");
			stateAustraliaSelect.prop('disabled', true).addClass("hidden");
			stateCanadaSelect.prop('disabled', false).removeClass("hidden");
		} else if(currentCountry == 'Australia') {
			flag = true;
			stateSelect.prop('disabled', true).addClass("hidden");
			stateUsaSelect.prop('disabled', true).addClass("hidden");
			stateCanadaSelect.prop('disabled', true).addClass("hidden");
			stateAustraliaSelect.prop('disabled', false).removeClass("hidden");
		} else {
			flag = false;
			stateCanadaSelect.prop('disabled', true).addClass("hidden");
			stateUsaSelect.prop('disabled', true).addClass("hidden");
			stateAustraliaSelect.prop('disabled', true).addClass("hidden");
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
			console.log(rmvFile);
			$.ajax({
				type: 'POST',
				url: '/businesses/deletefile',
				data: "id="+rmvFile+"&businessID="+$("#businessID").val()+"&imageType=gallery",
				dataType: 'html'
			}).done(function () { $(document).find(file.previewElement).remove(); $('.galleryContainer img[src*="'+rmvFile+'"]').closest(".galleryImageSIngle").remove(); loadBusinessGallery(); });
			//var _ref;
			
			//return (_ref = file.previewElement) != null ? _ref.parentNode.removeChild(file.previewElement) : void 0;   
		}
		console.log(this.options.maxFiles);
	}
	});
	
	$("#profileUpload").dropzone({
		url: "/businesses/addfiles",
		sending: function(file, xhr, formData){
			formData.append('imageType', "profile");
            formData.append('businessID', $("#businessID").val());
			formData.append('imageFormat', file.type);
			console.log(file.type);
        },
		addRemoveLinks : true,
		maxFiles:1,
		acceptedFiles: ".jpeg,.jpg,.png",
		init: function() {
			this.on('addedfile', function(file) {
				if (this.files.length > 1) {
				  this.removeFile(this.files[0]);
				  this.options.maxFiles = 1;
				}
			  });
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
	},
	removedfile: function(file) {
			var rmvFile = "";
		console.log("delete file length: "+fileList2.length);
		if(fileList2.length > 0) {
		for(f=0;f<fileList2.length;f++){

			if(fileList2[f].fileName == file.name)
			{
				console.log("new: "+ fileList2[f].fileName+"old: "+file.name);
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
			console.log(rmvFile);
			$.ajax({
				type: 'POST',
				url: '/businesses/deletefile',
				data: "id="+rmvFile+"&businessID="+$("#businessID").val()+"&imageType=profile",
				dataType: 'html'
			}).done(function () { $(document).find(file.previewElement).remove(); });
			//var _ref;
			
			//return (_ref = file.previewElement) != null ? _ref.parentNode.removeChild(file.previewElement) : void 0;   
		}
		console.log(this.options.maxFiles);
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
		maxFiles:1,
		acceptedFiles: ".jpeg,.jpg,.png",
		init: function() {
			this.on('addedfile', function(file) {
				if (this.files.length > 1) {
				  this.removeFile(this.files[0]);
				  this.options.maxFiles = 1;
				}
			  });
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
	},
	removedfile: function(file) {
			var rmvFile = "";
		console.log("delete file length: "+fileList3.length);
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
		}
	});
}

function loadBusinessGallery() {
	$(".businessGallery").fancybox({
		openEffect	: 'fade',
		closeEffect	: 'fade',
		type : "image",
		thumbs : {
			showOnStart : true,
			hideOnClosing : true
		}
	});
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