var bizID;
var userID;
var bName;
var dt;
var lastCountrySelected;
var textDescriptionInit = false;
var thumbnailUrls = [];
var coverThumbnailUrls = [];
var profileThumbnailUrls= [];
var businessModalFlag = false;
var editCommentFlag = false;
var socket = new Ws("ws://office.bunity.com:"+port+"/notifications");
var $userListTable = $('#userListTable');
var $activityTypeTable = $('#activityTypeTable');

$(document).ready( function() {
	if($userListTable.length) {
		UserListTable();
	}
	if($activityTypeTable.length) {
		settingsTable();
	}
	$('.selectAdmin').select2();
	if($("#flipbookStats").length) {
		highChartsAll();
	}
	//BusinessesListTable();
  // Warning alert
        

});

function UserListTable() {
	var showKey = "showBusinesses";
	var detailRows = [];
	
	
	//init table
	dt = $userListTable.DataTable( {
        "processing": true,
		"serverSide": true,
		"responsive": true,
		"keys": true,
		"autoWidth": false,
		"stateSave": false,
		"order": [[2, "asc"]],
		//"pagingType": "input",
        "ajax": "/userlist",
		"iDisplayLength": 2,
		"fixedHeader": true,
		"lengthMenu": [ [2, 5, 10, 25, 50, 100, -1], [2, 5, 10, 25, 50, 100, "All"] ],
		"rowId": "id",
		"select": {
          "style":    'os',
          "selector": 'tr>td:nth-child(1), tr>td:nth-child(3), tr>td:nth-child(4), tr>td:nth-child(5), tr>td:nth-child(6)'
        },
		
		/*'createdRow': function( row, data, dataIndex ) {
				$(row).attr('data-id', data.id);
		},*/
		"dom": '<"datatable-header"fl><"datatable-scroll-wrap"t><"datatable-footer"ip>',
        "language": {
            search: '<span>Filter:</span> _INPUT_',
            searchPlaceholder: 'Type to filter...',
            lengthMenu: '<span>Show:</span> _MENU_',
            paginate: { 'first': 'First', 'last': 'Last', 'next': '&rarr;', 'previous': '&larr;' }
        },
		"columns": [

			{
				"orderable": false,
                "data": "id",
				"defaultContent": "",
				"className": 'select-checkbox',
				"width": "90",
				'checkboxes': {
				    'selectRow': true
				}
			},
			{
                "class":          "details-control",
                "orderable":      false,
                "data":           null,
                "defaultContent": "",
				"width": "20"
            },
			{ "name": "firstname", "data": "firstname" },
			{ "name": "lastname", "data": "lastname" },
			{ "name": "email", "data": "email" },
			{ "name": "business", "defaultContent": "", "data": "Business.0.name", "orderable": true, "width": "auto"},
			{
                "class":          "actionBtns",
                "orderable":      false,
                "data":           null,
                "defaultContent": '<div class="btn-group actionBtnsContainer">'+
									'<button type="button" class="userEditBtn btn btn-sm btn-default"><i class="fa fa-pencil" aria-hidden="true"></i></button>' +
									'<button type="button" class="userDeleteBtn btn btn-sm btn-default"><i class="fa fa-trash" aria-hidden="true"></i></button>' +
									'</div>',
				"width": "140"
            },
		],
		"drawCallback": function () {
			dt.rows().every( function () {
				var rowS = this;
				var rowData = this.data();
				var hideShow = $(".hideShowBtn");
	
				if(rowData.Business.length < 1) {
					
					$userListTable.find("#"+rowData.id+" .details-control").removeClass("details-control");
				} else {
					//show by default all the businesses expanded
					rowS.child( format( rowData ) );
					if (sessionStorage.getItem(showKey) == "true") {
						$("#"+rowData.id).addClass( 'details' );
						hideShow.addClass("showBusiness");
						rowS.child.show();
					} else {
						rowS.child.hide();
					}
					
					
				}
			});
		}
        //"deferLoading": 57
    });

	dt.on( 'select deselect draw', function ( e, dt, type, indexes ) {
		console.log($("tbody .dt-checkboxes:checked").length);
		if ( $("tbody .dt-checkboxes:checked").length == $('tbody .dt-checkboxes').length) {
			$("thead tr").removeClass("indeterminate");
			$("thead tr").addClass("selected");
		} else {
			$("thead tr").removeClass("selected");
			if($("tbody .dt-checkboxes:checked").length > 0) {
				$("thead tr input").prop({
					indeterminate: true,
					checked: false
				});
				$("thead tr").addClass("indeterminate");
			} else {
				$("thead tr").removeClass("selected");
				$("thead tr").removeClass("indeterminate");
			}
		}
		
	});
	
	//expand business button
	$userListTable.on( 'click', 'tr td.details-control', function () {
        var tr = $(this).closest('tr');
        var row = dt.row( tr );
		
        var idx = $.inArray( tr.attr('id'), detailRows );
 
        if ( row.child.isShown() ) {
            tr.removeClass( 'details' );
            row.child.hide();
 
            // Remove from the 'open' array
            detailRows.splice( idx, 1 );
        }
        else {
            tr.addClass( 'details' );
            row.child( format( row.data() ) ).show();
 
            // Add to the 'open' array
            if ( idx === -1 ) {
                detailRows.push( tr.attr('id') );
            }
        }
    });
	
	//trigger click when chainging view
	dt.on( 'draw', function () {
        $.each( detailRows, function ( i, id ) {
			setTimeout(function() {
				 $('#'+id+' td.details-control').trigger( 'click' );
			},50)
           
        } );
    } );
	
	
	//deleteBusiness
	$("#userListTable").on("click", ".businessDeleteBtn", function() {
		userID = $(this).parents('tr').last().prev().attr("id");
		bizID = $(this).closest('tr').data("id");
		console.log(bizID);
		//$('.popupContainer').html(confirmDeleteModal);
		bName = $(this).closest('tr').find(".tdBusinessName").text();
		$(".deleteBusinessName strong").text(bName);
		
		//$('.confirmDeleteModal').modal();
		confirmBusinessDelete(bName);
			
	});
	
	
	
	//hide modal event
	$(".popupContainer").on('hidden.bs.modal', '.modal', function (e) {
		businessModalFlag = false;
		$(this).remove();
		$(".dz-hidden-input").remove();
		tinymce.execCommand('mceRemoveEditor',true,'businessDescription');
	})
	
	//business edit action
	$("#userListTable").on("click", ".businessEditBtn", function() {
		//$(".businessEditModal").remove();
		var btnEdit = $(this);
		businessModalFlag = true;
		bName = $(this).closest('tr').find(".tdBusinessName").text();
		btnEdit.attr("disabled", true).prepend('<i class="icon-spinner2 spinner position-left"></i>');
		$('.popupContainer').html(businessEditModal);
		$('.businessCountry').select2();
		bizID = $(this).closest('tr').data("id");
		userID = $(this).parents('tr').last().prev().attr("id");
		$.ajax({
			type:"GET",
			url:"/business/"+bizID,
			success: function(data) {	
				var b = data.business;
				var statesInput;
				var flag = false;
				if(statesUSA.includes(b.state) && (b.country == "Australia" || b.country == "United States" || b.country == "Canada")) {
					statesInput = StatesSelectHTML(statesUSA);
					flag = true;
				} else if(statesAustralia.includes(b.state) && (b.country == "Australia" || b.country == "United States" || b.country == "Canada")) {
					statesInput = StatesSelectHTML(statesAustralia);
					flag = true;
				} else if(statesCanada.includes(b.state) && (b.country == "Australia" || b.country == "United States" || b.country == "Canada")) {
					statesInput = StatesSelectHTML(statesCanada);
					flag = true;
				} else {
					statesInput = StatesSelectHTML(false);
				}
				
				lastCountrySelected = b.country;
				$('.businessStateContainer').append(statesInput);
				if(flag) {
					$('.businessState').select2();
				}
				
				$(".businessCountry").val(b.country).trigger('change.select2');
				$(".businessState").val(b.state).trigger('change.select2');

				
				$('.businesscateg, .businesscateg2').select2();
				if (in_array(businessCategs, b.category)) {
					$(".businesscateg").val(b.category).trigger('change.select2');
				}
				if (in_array(businessCategs, b.category2)) {
					$(".businesscateg2").val(b.category2).trigger('change.select2');
				}
				
				if(yearsBusiness.includes(b.yearsBusiness)) {
					$(".yearsBusiness").val(b.yearsBusiness);
				}
				
				if(numberEmployees.includes(b.numberEmployees)) {
					$(".numberEmployees").val(b.numberEmployees);
				}
				
				if(sizeBusiness.includes(b.sizeBusiness)) {
					$(".sizeBusiness").val(b.sizeBusiness);
				}
				
				if(relationshipBusiness.includes(b.relationshipBusiness)) {
					$(".relationshipBusiness").val(b.relationshipBusiness);
				}
				

				$("#businessDescription").val(b.description);
				initDescriptionBox();
				
				
				$("#businessEditModalLabel strong").text(b.name);
				$(".businessName").val(b.name);
				$(".businessPhone").val(b.phone);
				$(".businessCountry").val(b.country);
				$(".businessState").val(b.state);
				$(".businessArea").val(b.area);
				$(".businessCity").val(b.city);
				$(".businessAddress").val(b.address);
				$(".businessAddress2").val(b.address2);
				$(".businessWebsite").val(b.website);
				$(".businessPostalCode").val(b.postalcode);
				
				$(".businessFacebook").val(b.Facebook);
				$(".businessGoogle").val(b.google);
				$(".businessInstagram").val(b.instagram);
				$(".businessYoutube").val(b.youtube);
				$(".businessPinterest").val(b.pinterest);
				$(".businessLinkedin").val(b.linkedin);
				$(".businessTwitter").val(b.twitter);
				
				thumbnailUrls = b.gallery;
				coverThumbnailUrls = b.cover;
				profileThumbnailUrls = b.profile;
				galleryUploadFunc(bizID, userID);
				
				btnEdit.attr("disabled", false);
				btnEdit.find(".icon-spinner2").remove();
				//textarea comment
				//$(".chatInput").autogrow();
				showBusinessComments(data.comments);
				$(".activityTypeSelect").select2();
				//#atjs
				$('#atjs').atwho({
					at: "#",
					data: activitesType,
					startWithSpace: false,
					displayTpl: "<li> ${name} </li>",
					insertTpl: '<span data-id="${id}" class="tag label label-primary activityTypeLabel"> ${name}</span></div>',
					limit: 10000,
					callbacks: {
					  filter: function(query, data, searchKey)
					  {
	
						if($('#atjs .activityTypeLabel').length  < 1) {
							return data;
						}
					  } 
					}
				});
				$('.businessEditModal').modal();
			},
			dataType: 'json',
		  });
		  
		  
		  
	})
	
	//submit business chat
	$('.popupContainer').on("keydown", "#atjs", function(event) {
		var msg = $("#atjs").text();
		var msgHTML = $("#atjs").html()
		var tempMsg = msg;
		if (event.keyCode == 13 && event.shiftKey) {

		} else if (event.keyCode == 13) {
			var parentID = getCommentID(msg);
			if(parentID) {
				var msgReply = msg.replace("@"+parentID, "");
				var msgReplyHTML = msgHTML.replace("@"+parentID, "");
				
				if(msgReply.replace(/\s/g, '').length) {
					submitChat(msgReplyHTML.slice(1), parentID);
					
				}
				$(this).html("");
				return false;
			}
			
			if(isEditActive(tempMsg)) {
				var msgEdit = msg.replace("~", "");
				var msgEditHTML = msgHTML.replace("~", "");
				if (msgEdit.replace(/\s/g, '').length) {
					updateComm(msgEditHTML);
				}
				$(this).html("");
				$(".media-list .media.edit").removeClass("edit");
				editCommentFlag = false;
				return false;
				
			}
			if (tempMsg.replace(/\s/g, '').length) {
				submitChat(msgHTML);
			} else {
				$(this).html("");
			}
			return false;
		}
		
	});

	
	//business update event
	$(".popupContainer").on("click", ".businessUpdate", function() {
		var btnUpdate = $(this);
		
		console.log(businessModalFlag);
		btnUpdate.attr("disabled", true).prepend('<i class="icon-spinner2 spinner position-left"></i>');
		var formSelector = $(".editBusinessForm");
		var b = {};
		$.each($(formSelector).serializeArray(), function(i, field) {
			b[field.name] = field.value;
		});
		b.businessdescription = tinymce.activeEditor.getContent();
		b.userid = userID;
		$.ajax({
		type:"POST",
		url:"/business/"+bizID,
			data: {b},
			success: function(data) {
				btnUpdate.attr("disabled", false)
				btnUpdate.find(".icon-spinner2").remove();
				$(".businessEditModal").modal("hide")
				dt.ajax.reload(null, false);
				var opts = {};
				opts.title = "Busienss Updated";
				opts.text = "Business <strong>"+bName+"</strong> was successfully updated!";
				opts.type = "success";
				opts.icon = 'icon-checkmark3';
				new PNotify(opts);
			},
			dataType: 'json',
		  });
	});
	
	//hide/show businesses
	$(".hideShowBtn").on("click", function(e) {
		e.preventDefault();
		var showBusinesses;
		if (typeof(Storage) !== "undefined") {
			if (sessionStorage.getItem(showKey) == "true") {
				sessionStorage.setItem(showKey, false);
			} else {
				sessionStorage.setItem(showKey, true);
			}
			showBusinesses = sessionStorage.getItem(showKey);
			var hideShow = $(".hideShowBtn");
			dt.rows().every( function () {
				var row = this;
				var rowData = this.data();
				if(rowData.Business.length > 0) {
					if(showBusinesses == "true") {
						row.child.show();
						hideShow.addClass("showBusiness");
						$("#"+rowData.id).addClass( 'details' );
					} else {
						row.child.hide();
						hideShow.removeClass("showBusiness");

						$("#"+rowData.id).removeClass( 'details' );
					}
				}
			});
		}
	});
	
	$(".popupContainer").on("select2:select", ".businessCountry", function(){
		var selectedCountry = $(this).val();
		var flag = false;
		var statesInput = "";
		if(selectedCountry == 'United States'){
			statesInput = StatesSelectHTML(statesUSA);
			flag = true;
		} else if(selectedCountry == 'Canada') {
			statesInput = StatesSelectHTML(statesCanada);
			flag = true;
		} else if(selectedCountry == 'Australia') {
			statesInput = StatesSelectHTML(statesAustralia);
			flag = true;
		} else {
			statesInput = StatesSelectHTML(false);
		}
		
		if (flag) {
			
			$(".businessState, .businessStateContainer .select2").remove();
			$('.businessStateContainer').append(statesInput);
			$('.businessState').select2();
		}
		if (!flag && (lastCountrySelected == "Australia" || lastCountrySelected == "United States" || lastCountrySelected == "Canada")) {
			$(".businessState, .businessStateContainer .select2").remove();
			$('.businessStateContainer').append(statesInput);
			
		}
		lastCountrySelected = selectedCountry;
	});
	
	$(document).on("click", ".replyBtn", function() {
		commentReply($(this));
	});
	
	$(document).on("click", ".editCommBtn", function() {
		editComm($(this));
	});
	
	
	
	socket.On("newBusinessChat", function (comment) {
		var comment = JSON.parse(comment);
		var notice = new PNotify({
			title: 'Message reply',
			text: comment.text,
			icon: 'icon-comment-discussion'
		});
		notice.get().click(function() {
			notice.remove();
		});  
	});
	
	socket.On("refreshCommBiz", function (comments) {
		if(businessModalFlag) {
			var comments = JSON.parse(comments);
			if(comments[0].business_id === bizID) {
				showBusinessComments(comments);
			}
		}
	});
	
	$(document).on("input", "#atjs", function() {
		var _this = $(this);
		if(editCommentFlag) {
			var val = _this.text();
			if(val.charAt(0) != "~") {
				editCommentFlag = false;
				$(".media-list .media.edit").removeClass("edit");
			}
		}
	});
		
}



//business 
var businessEditModal = '<div class="modal fade businessEditModal ">'+
						'<div class="modal-dialog modal-lg">'+
							'<div class="modal-content">'+
								'<div class="modal-header">'+
									'<button type="button" class="close" data-dismiss="modal">×</button>'+
									'<h5 id="businessEditModalLabel" class="modal-title">Update business - <strong></strong></h5>'+
								'</div>'+
								
								'<div class="modal-body">'+
									'<div class="row">'+
									'<div class="col-sm-8">'+
									'<ul class="nav nav-tabs">'+
										'<li class="active"><a href="#nameAddressStepTab" data-toggle="tab">Name & Address</a></li>'+
										'<li><a href="#detailsTab" data-toggle="tab">Details</a></li>'+
										'<li><a href="#description" data-toggle="tab">Description</a></li>'+
										'<li><a href="#social" data-toggle="tab">Social</a></li>'+
										'<li><a href="#images" data-toggle="tab">Images</a></li>'+
									'</ul>'+
									'<form class="editBusinessForm">' +
										'<div class="tab-content">'+
											'<div class="tab-pane active" id="nameAddressStepTab">'+
												'<div class="form-group">'+
													'<div class="row">'+
														'<div class="col-sm-6">'+
															'<label>Name</label>'+
															'<input type="text" placeholder="" name="name" class="form-control businessName">'+
														'</div>'+

														'<div class="col-sm-6">'+
															'<label>Phone #</label>'+
															'<input type="text" placeholder="" name="phone" class="form-control businessPhone">'+
														'</div>'+
													'</div>'+
												'</div>'+

												'<div class="form-group">'+
													'<div class="row">'+
														'<div class="col-sm-6">'+
															'<label>Address line 1</label>'+
															'<input type="text" placeholder="" name="address" class="form-control businessAddress">'+
														'</div>'+

														'<div class="col-sm-6">'+
															'<label>Address line 2</label>'+
															'<input type="text" placeholder="" name="address2" class="form-control businessAddress2">'+
														'</div>'+
													'</div>'+
												'</div>'+
												
												'<div class="form-group">'+
													'<div class="row">'+
														'<div class="col-sm-4">'+
															'<label>Country</label>'+
															CountriesSelectHTML()+
														'</div>'+

														'<div class="col-sm-4 businessStateContainer">'+
															'<label>State/Area</label>'+
															
														'</div>'+

														'<div class="col-sm-4">'+
															'<label>City</label>'+
															'<input type="text" placeholder="" name="city" class="form-control businessCity">'+
														'</div>'+
													'</div>'+
												'</div>'+

												'<div class="form-group">'+
													'<div class="row">'+
														'<div class="col-sm-6">'+
															'<label>ZIP/Postal code</label>'+
															'<input type="text" placeholder="" name="postal_code" class="form-control businessPostalCode">'+
														'</div>'+

														'<div class="col-sm-6">'+
															'<label>Website</label>'+
															'<input type="text" placeholder="" name="website" class="form-control businessWebsite">'+
														'</div>'+
													'</div>'+
												'</div>'+
											'</div>'+

											'<div class="tab-pane" id="detailsTab">'+
												'<div class="form-group">'+
													'<div class="row">'+
														'<div class="col-sm-6">'+
															'<label>Primary Category</label>'+
															BusinessCategsHTML("businesscateg")+
														'</div>'+
														'<div class="col-sm-6">'+
															'<label>Secondary Category</label>'+
															BusinessCategsHTML("businesscateg2")+
														'</div>'+
													'</div>'+
												'</div>'+
													
												'<div class="form-group">'+
													'<div class="row">'+
														'<div class="col-sm-6">'+
															'<label>Years in Business</label>'+
															OtherSelectsHTML(yearsBusiness, "yearsBusiness", "Years in Business")+
														'</div>'+
														'<div class="col-sm-6">'+
															'<label>Number of Employees</label>'+
															OtherSelectsHTML(numberEmployees, "numberEmployees", "Number of Employees")+
														'</div>'+
													'</div>'+
												'</div>'+
												
												'<div class="form-group">'+
													'<div class="row">'+
														'<div class="col-sm-6">'+
															'<label>Size of Business</label>'+
															OtherSelectsHTML(sizeBusiness, "sizeBusiness", "Years in Business")+
														'</div>'+
														'<div class="col-sm-6">'+
															'<label>Secondary Category</label>'+
															OtherSelectsHTML(relationshipBusiness, "relationshipBusiness", "What is your relationship to this business?")+
														'</div>'+
													'</div>'+
												'</div>'+
											'</div>'+
											
											'<div class="tab-pane" id="description">'+
												'<div class="form-group">'+
													'<textarea id="businessDescription" name="businessdescription" class="form-control"></textarea>'+
												'</div>'+
											'</div>'+
											
											'<div class="tab-pane" id="social">'+
												'<div class="form-group">'+
													'<div class="row">'+
														'<div class="col-sm-6">'+
															'<label>Facebook</label>'+
															'<input type="text" placeholder="" name="businessFacebook" class="form-control businessFacebook">'+
														'</div>'+

														'<div class="col-sm-6">'+
															'<label>Google+</label>'+
															'<input type="text" placeholder="" name="businessGoogle" class="form-control businessGoogle">'+
														'</div>'+
													'</div>'+
												'</div>'+
											
												'<div class="form-group">'+
													'<div class="row">'+
														'<div class="col-sm-6">'+
															'<label>Instagram</label>'+
															'<input type="text" placeholder="" name="businessInstagram" class="form-control businessInstagram">'+
														'</div>'+

														'<div class="col-sm-6">'+
															'<label>YouTube</label>'+
															'<input type="text" placeholder="" name="businessYoutube" class="form-control businessYoutube">'+
														'</div>'+
													'</div>'+
												'</div>'+
											
												'<div class="form-group">'+
													'<div class="row">'+
														'<div class="col-sm-6">'+
															'<label>Pinterest</label>'+
															'<input type="text" placeholder="" name="businessPinterest" class="form-control businessPinterest">'+
														'</div>'+

														'<div class="col-sm-6">'+
															'<label>Linkedin</label>'+
															'<input type="text" placeholder="" name="businessLinkedin" class="form-control businessLinkedin">'+
														'</div>'+
													'</div>'+
												'</div>'+
											
												'<div class="form-group">'+
													'<div class="row">'+
														'<div class="col-sm-6">'+
															'<label>Twitter</label>'+
															'<input type="text" placeholder="" name="businessTwitter" class="form-control businessTwitter">'+
														'</div>'+

														'<div class="col-sm-6">'+
															
														'</div>'+
													'</div>'+
												'</div>'+
											'</div>'+
											
											'<div class="tab-pane" id="images">'+
												'<div class="editPhotoTitle">Gallery (max 8 files):</div>'+
												'<div id="galleryUpload" class="dropzone"></div>'+
												'<div class="editPhotoTitle">Profile (max 1 file):</div>'+
												'<div id="profileUpload" class="dropzone"></div>'+
												'<div class="editPhotoTitle">Cover (max 1 file):</div>'+
												'<div id="coverUpload" class="dropzone"></div>'+
											'</div>'+		
										'</div>'+
									'</form>'+
									'</div>'+
									'<div class="col-sm-4">'+
										//'<textarea name="chat" class="form-control content-group chatInput" placeholder="Enter your message..."></textarea>'+
										'<div id="atjs" class="inputor form-control" contenteditable="true"></div>'+
										//OtherSelectsHTML(activitesType, "activityTypeSelect", "activity type")+
										'<ul class="media-list chat-list content-group bMessageContainer"></ul>'+
									'</div>'+
									'</div>'+
								'</div>'+
								'<div class="modal-footer">'+
									'<button type="button" class="btn btn-link" data-dismiss="modal">Close</button>'+
									'<button type="submit" class="btn btn-primary businessUpdate">Update</button>'+
								'</form>'+
							'</div>'+
						'</div>'+
					'</div>'+
				'</div>';

var confirmDeleteModal = '<div class="modal fade confirmDeleteModal" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel" aria-hidden="true">'+
	 '<div class="modal-dialog modal-md">'+
		'<div class="modal-content">'+
			'<div class="modal-header">'+
				'<h5 class="modal-title" id="confirmDeleteModalLabel">Confirm Delete <span></span></h5>'+
			'</div>'+
			'<div class="modal-body">'+
				'<p class="deleteBusinessName">You are about to delete a business: <strong></strong>, this procedure is irreversible.</p>'+
				'<p>Do you want to proceed?</p>'+
			'</div>'+
			'<div class="modal-footer">'+
				'<button type="button" class="btn btn-secondary" data-dismiss="modal">Cancel</button>'+
				'<button type="button" class="btn btn-danger confirmBDelete">Delete</button>'+
			'</div>'+
		'</div>'+
	  '</div>'+
	'</div>';
	
function format ( d ) {
	var b = d.Business;
	var html = '<table class="businessTable"><tbody>';
	for (var i=0; i< b.length; i++) {
		html += '<tr data-id="'+b[i].id+'"><td class="tdBusinessName">'+b[i].name+'</td><td class="actionBtns"><div class="btn-group btn-group-sm actionBtnsContainer">'+
									'<button type="button" class="businessEditBtn btn btn-sm btn-default"><i class="fa fa-pencil" aria-hidden="true"></i></button>' +
									'<button type="button" class="businessDeleteBtn btn btn-sm btn-default"><i class="fa fa-trash" aria-hidden="true"></i></button>' +		
									'</div></td></tr>';
	}
	html += '</tbody></table>';
    return html;
}

//delete business callback
function DeleteBusiness() {
	$(".sweet-alert .confirm").attr("disabled", true);
	$.ajax({
		type:"DELETE",
		url:"/business/"+bizID, //+"/"+bizID,
		success: function(data) {
			if (data.success === true) {
				//$(".sweet-alert .confirm").attr("disabled", false);
				swal.close();
				dt.ajax.reload(null, false);
				var opts = {};
				opts.title = "Busienss Deleted";
				opts.text = "Business <strong>"+bName+"</strong> was successfully deleted!";
				opts.type = "info";
				opts.icon = 'icon-trash';
				new PNotify(opts);
				//$('.confirmDeleteModal').modal("hide");
			}
		},
		dataType: 'json',
	});
}


function confirmBusinessDelete(bName) {
	swal({
		title: "Are you sure?",
		text: "You will not be able to recover <strong>"+bName+"</strong> business!",
		html: true,
		type: "warning",
		showCancelButton: true,
		closeOnConfirm: false,
		confirmButtonColor: "#FF7043",
		confirmButtonText: "Yes, delete it!"
	},
	function(isConfirm){
		if (isConfirm) {
			DeleteBusiness();
		}
		else {
			
		}
	});
}

function CountriesSelectHTML() {
	var html = '<select class="form-control businessCountry" name="country"><option value="">Choose Country</option>';
	for(var i=0; i < countryList.length; i++) {
		var country = countryList[i];
		html += '<option value="'+country+'">'+country+'</option>';
	}
	html += '</select>'
	return html;
}

function StatesSelectHTML(states) {
	if (Array.isArray(states)) {
		var html = '<select class="form-control businessState" name="state"><option value="">Choose State</option>';
		for(var i=0; i < states.length; i++) {
			var state = states[i];
			html += '<option value="'+state+'">'+state+'</option>';
		}
		html += '</select>'
	} else {
		html = '<input type="text" class="form-control businessState" name="state">';
	}
	return html;
}

function BusinessCategsHTML(name) {
	var html = '<select class="form-control '+name+'" name="'+name+'"><option value="">Choose Category</option>';
	for(var i=0; i < businessCategs.length; i++) {
		var categ = businessCategs[i].name;
		html += '<option value="'+categ+'">'+categ+'</option>';
	}
	html += '</select>'
	return html;
}

function OtherSelectsHTML(array, name, dummyOption) {
	var html = '<select class="form-control '+name+'" name="'+name+'"><option value="">Choose '+dummyOption+'</option>';
	for(var i=0; i < array.length; i++) {
		var item = array[i];
		if (name == 'activityTypeSelect') {
			html += '<option value="'+item.name+'">'+item.name+'</option>';
		} else {
			html += '<option value="'+item+'">'+item+'</option>';
		}
	}
	html += '</select>'
	return html;
}


function in_array(array, name) {
    for(var i=0;i<array.length;i++) {
		if (array[i].name == name) {
			return true;
		}
    }
    return false;
}

function initDescriptionBox() {
	if (textDescriptionInit) {
		tinymce.execCommand('mceAddEditor',true,'businessDescription');
		console.log(1);
	} else {
	//if($("#businessDescription").length > 0) {
		tinymce.init({
			selector: '#businessDescription',
			branding: false,
			height: 300,
			paste_as_text: true,
			menubar: false,
			plugins: [
				'paste lists',
				'wordcount'
			],
			setup: function(editor) {
			editor.on('init', function() {
				textDescriptionInit = true;
			});
		  },
			toolbar: 'bold italic | bullist numlist',
		});
	}
		
	//}
}

function galleryUploadFunc(bizID, userID) {
	var cachedFilename;
	
	Dropzone.autoDiscover = false;
	var fileList = new Array;
	var fileList2 = new Array;
	var fileList3 = new Array;
	$("#galleryUpload").dropzone({
		url: "/picture/add/"+bizID,
		sending: function(file, xhr, formData){
			formData.append('imageType', "gallery");
			formData.append('imageFormat', file.type);
			formData.append('userID', userID);
        },
		/*headers: {
			'Cache-Control': null,
			'X-Requested-With': null
		},*/
		addRemoveLinks : true,
		maxFiles:8,
		acceptedFiles: ".jpeg,.jpg,.png",
		init: function() {
			var myDropzone = this;
			var existingFileCount = thumbnailUrls.length;
			//myDropzone.options.maxFiles = myDropzone.options.maxFiles - existingFileCount ;
			if (thumbnailUrls) {
				for (var i = 0; i < thumbnailUrls.length; i++) {
					var imgURL = domain+"/static/uploads/"+userID+"/"+bizID+"/gallery/"+thumbnailUrls[i];
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
		var opts = {};
		opts.title = "Error uploading gallery image";
		opts.text = message;
		opts.type = "error";
		opts.icon = 'icon-blocked';
		new PNotify(opts);
		this.removeFile(file);
		$(document).find(file.previewElement).remove(); 
	},
	processing:function(file) {
		file.newName = '_' + Math.random().toString(36).substr(2, 9) +"."+ file.name.split('.').pop();
		if (file.previewElement) {
          file.previewElement.classList.add("dz-processing");
          if (file._removeLink) {
            return file._removeLink.textContent = this.options.dictCancelUpload;
          }
        }
	},
	
	maxfilesexceeded: function(file) {
		$(document).find(file.previewElement).remove(); 
	},
	success: function(file, serverFileName) {
		fileList.push ({"serverFileName" : serverFileName.fname, "fileName" : file.newName});
		var opts = {};
		opts.title = "Image Uploaded!";
		opts.text = "The gallery image was updated successfully";
		opts.type = "success";
		opts.icon = 'icon-checkmark3';
		new PNotify(opts);
	},

	removedfile: function(file) {
		var rmvFile = "";
		console.log("delete file length: "+fileList.length);
		if(fileList.length > 0) {
		for(f=0;f<fileList.length;f++){

			if(fileList[f].fileName == file.newName)
			{
				console.log("new: "+ fileList[f].fileName+"old: "+file.name);
				rmvFile = fileList[f].serverFileName;
				fileList.splice(f,1);
				//myDropzone.options.maxFiles = myDropzone.options.maxFiles + 1;
				break;
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
				url: '/picture/delete',
				data: "businessID="+bizID+"&userID="+userID+"&fileID="+rmvFile+"&imageType=gallery",
				dataType: 'json'
			}).done(
			function (data) {
				if(data.success) {
					$(document).find(file.previewElement).remove(); 
					$('.galleryContainer img[src*="'+rmvFile+'"]').closest(".galleryImageSIngle").remove();
					var opts = {};
					opts.title = "Image removed!";
					opts.text = "The gallery image was removed successfully";
					opts.type = "info";
					opts.icon = 'icon-checkmark3';
					new PNotify(opts);
				}
			});
			//var _ref;
			
			//return (_ref = file.previewElement) != null ? _ref.parentNode.removeChild(file.previewElement) : void 0;   
		}
		//console.log(this.options.maxFiles);
	}
	});
	
	
	/* PROFILE UPLOAD */
	$("#profileUpload").dropzone({ 
		url: "/picture/add/"+bizID,
		sending: function(file, xhr, formData){
			formData.append('imageType', "profile");
			formData.append('imageFormat', file.type);
			formData.append('userID', userID);
        },
		addRemoveLinks : true,
		maxFiles:1,
		maxFilesize: 5,
		acceptedFiles: ".jpeg,.jpg,.png",
		
		accept: function(file, done) {
		if (file.width < 160 && file.height < 160) {
			file.rejectDimensions("The resolution has to be at least 160 x 160");
			}	else {
				done();
			}
		},
		init: function() {			
			var myDropzone = this;
			if (profileThumbnailUrls) {
				for (var i = 0; i < profileThumbnailUrls.length; i++) {
					var imgURL = domain+"/static/uploads/"+userID+"/"+bizID+"/profile/"+profileThumbnailUrls[i];
					var mockFile = { 
						name: profileThumbnailUrls[i], 
						//size: 12345, 
						//type: 'image/jpeg', 
						status: Dropzone.ADDED, 
						server:true,
						accepted: true,
						url: imgURL,
					};

					// Call the default addedfile event handler
					myDropzone.emit("addedfile", mockFile);
					myDropzone.emit("complete", mockFile);
					// And optionally show the thumbnail of the file:
					myDropzone.emit("thumbnail", mockFile, imgURL);

					myDropzone.files.push(mockFile);


				}
			}
		},
	success: function(file, serverFileName) {
		fileList2.push ({"serverFileName" : serverFileName.fname, "fileName" : file.newName});
	
		if(serverFileName.fname) {
			var opts = {};
			opts.title = "Image Uploaded!";
			opts.text = "The profile image was updated successfully";
			opts.type = "success";
			opts.icon = 'icon-checkmark3';
			new PNotify(opts);
		}
	},
	
	error: function(file, message, xhr) {
		var opts = {};
		opts.title = "Error uploading profile image";
		opts.text = message;
		opts.type = "error";
		opts.icon = 'icon-blocked';
		new PNotify(opts);
		this.removeFile(file);
		$(document).find(file.previewElement).remove(); 
	},
	processing:function(file) {
		file.newName = '_' + Math.random().toString(36).substr(2, 9) +"."+ file.name.split('.').pop();
		
		if (file.previewElement) {
          file.previewElement.classList.add("dz-processing");
          if (file._removeLink) {
            return file._removeLink.textContent = this.options.dictCancelUpload;
          }
        }
	},
	maxfilesexceeded: function(file) {
		this.removeFile(file);
		$(document).find(file.previewElement).remove(); 
	},
	removedfile: function(file) {
		var myDropzone = this
		var rmvFile = "";
		if(fileList2.length > 0) {
			for(f=0;f<fileList2.length;f++){

				if(fileList2[f].fileName == file.newName)
				{
					console.log("new: "+ fileList2[f].fileName+"old: "+file.name);
					rmvFile = fileList2[f].serverFileName;
					fileList2.splice(f,1);
					//myDropzone.options.maxFiles = myDropzone.options.maxFiles + 1;
					break;
				}

			}
		}
		if(rmvFile == "" && file.server) {
			rmvFile = file.name;
		}
		if (rmvFile && file.accepted){
			
			//console.log(rmvFile);
			$.ajax({
				type: 'POST',
				url: '/picture/delete',
				data: "businessID="+bizID+"&userID="+userID+"&fileID="+rmvFile+"&imageType=profile",
				dataType: 'json'
			}).done(
			function (data) {
				if(data.success) {
					$(document).find(file.previewElement).remove(); 
					var opts = {};
					opts.title = "Image removed!";
					opts.text = "The profile image was removed successfully";
					opts.type = "info";
					opts.icon = 'icon-checkmark3';
					new PNotify(opts);
				}
			});
			//var _ref;
			
			//return (_ref = file.previewElement) != null ? _ref.parentNode.removeChild(file.previewElement) : void 0;   
		} else {
			$(document).find(file.previewElement).remove();
		}
		//console.log(this.options.maxFiles);
	}

	});
	
	


	/* COVER UPLOAD */
	$("#coverUpload").dropzone({ 
		url: "/picture/add/"+bizID,
		sending: function(file, xhr, formData){
			formData.append('imageType', "cover");
			formData.append('imageFormat', file.type);
			formData.append('userID', userID);
        },
		addRemoveLinks : true,
		maxFiles:1,
		maxFilesize: 5,
		acceptedFiles: ".jpeg,.jpg,.png",
		
		accept: function(file, done) {
		if (file.width < 840 && file.height < 285) {
			file.rejectDimensions("The resolution has to be at least 840 x 285");
			}	else {
				done();
			}
		},
		init: function() {			
			var myDropzone = this;
			if (coverThumbnailUrls) {
				for (var i = 0; i < coverThumbnailUrls.length; i++) {
					var imgURL = domain+"/static/uploads/"+userID+"/"+bizID+"/cover/"+coverThumbnailUrls[i];
					var mockFile = { 
						name: coverThumbnailUrls[i], 
						//size: 12345, 
						//type: 'image/jpeg', 
						status: Dropzone.ADDED, 
						server:true,
						accepted: true,
						url: imgURL,
					};

					// Call the default addedfile event handler
					myDropzone.emit("addedfile", mockFile);
					myDropzone.emit("complete", mockFile);
					// And optionally show the thumbnail of the file:
					myDropzone.emit("thumbnail", mockFile, imgURL);

					myDropzone.files.push(mockFile);


				}
			}
		},
	success: function(file, serverFileName) {
		fileList3.push ({"serverFileName" : serverFileName.fname, "fileName" : file.newName});
	
		if(serverFileName.fname) {
			var opts = {};
			opts.title = "Image Uploaded!";
			opts.text = "The cover image was updated successfully";
			opts.type = "success";
			opts.icon = 'icon-checkmark3';
			new PNotify(opts);
		}
	},
	
	error: function(file, message, xhr) {
		var opts = {};
		opts.title = "Error uploading cover image";
		opts.text = message;
		opts.type = "error";
		opts.icon = 'icon-blocked';
		new PNotify(opts);
		this.removeFile(file);
		$(document).find(file.previewElement).remove(); 
	},
	processing:function(file) {
		file.newName = '_' + Math.random().toString(36).substr(2, 9) +"."+ file.name.split('.').pop();
		
		if (file.previewElement) {
          file.previewElement.classList.add("dz-processing");
          if (file._removeLink) {
            return file._removeLink.textContent = this.options.dictCancelUpload;
          }
        }
	},
	maxfilesexceeded: function(file) {
		this.removeFile(file);
		$(document).find(file.previewElement).remove(); 
	},
	removedfile: function(file) {
		var myDropzone = this
		var rmvFile = "";
		if(fileList3.length > 0) {
			for(f=0;f<fileList3.length;f++){

				if(fileList3[f].fileName == file.newName)
				{
					console.log("new: "+ fileList3[f].fileName+"old: "+file.name);
					rmvFile = fileList3[f].serverFileName;
					fileList3.splice(f,1);
					//myDropzone.options.maxFiles = myDropzone.options.maxFiles + 1;
					break;
				}

			}
		}
		if(rmvFile == "" && file.server) {
			rmvFile = file.name;
		}
		if (rmvFile && file.accepted){
			
			//console.log(rmvFile);
			$.ajax({
				type: 'POST',
				url: '/picture/delete',
				data: "businessID="+bizID+"&userID="+userID+"&fileID="+rmvFile+"&imageType=cover",
				dataType: 'json'
			}).done(
			function (data) {
				if(data.success) {
					$(document).find(file.previewElement).remove(); 
					var opts = {};
					opts.title = "Image removed!";
					opts.text = "The cover image was removed successfully";
					opts.type = "info";
					opts.icon = 'icon-checkmark3';
					new PNotify(opts);
				}
			});
			//var _ref;
			
			//return (_ref = file.previewElement) != null ? _ref.parentNode.removeChild(file.previewElement) : void 0;   
		} else {
			$(document).find(file.previewElement).remove();
		}
		//console.log(this.options.maxFiles);
	}

	});
}

function submitChat(msg, parentID) {
	
	var data = {
		"bizID": bizID,
	}
	if (parentID != "") {
		data.parentID = parentID;
		msg = msg.replace("@"+parentID, "");
	}
	if($("#atjs .activityTypeLabel").length > 0) {
		data.activityTypeID = $("#atjs .activityTypeLabel").last().data('id');
		data.activityTypeName = $("#atjs .activityTypeLabel").last().text().trim();
	}
	$("#atjs").text("");
	data.msg = msg;
	var nowTime = timeSince(new Date())
	
	$.ajax({
		type:"POST",
		url:"/business/comment",
		data: data,
		success: function(data) {
			if(data.success) {
				if(data.data.parent_id) {
					$("li[data-id='"+data.data.parent_id+"']").after(chatMsgTemplate(msg, userName, "reversed", nowTime, data.data.id, "child"));
				} else {
					$(".bMessageContainer").prepend(chatMsgTemplate(msg, userName, "reversed", nowTime, data.data.id, "parent"));
				}
				socket.OnConnect(function () {
					socket.Emit("newBusinessChat", data.data.id);
					socket.Emit("refreshCommBiz", bizID);
				});
				
			}
		},
		dataType: 'json',
	});
}

function showBusinessComments(comments) {
	var commentType;
	if(comments.length) {
		$(".bMessageContainer").html("");
		for(var i=0; i < comments.length; i++) {
			var comment = comments[i];
			var date = new Date(comment.time);
			var sinceTime = timeSince(date);
			commentType = "";
			if (userid == comment.author.id ) {
				commentType = "reversed";
			}
			if(comment.parent_id) {
				$("li[data-id='"+comment.parent_id+"']").after(chatMsgTemplate(comment.text, comment.author.name, commentType, sinceTime, comment.id, "child"));
			} else {
				$(".bMessageContainer").prepend(chatMsgTemplate(comment.text, comment.author.name, commentType, sinceTime, comment.id, "parent"));
			}
		}
	}
}

$.fn.extend({
	autogrow: function () {
		$(this).on('change keyup keydown paste cut', function () {
			$(this).height(0).height(this.scrollHeight - 10);
		}).change();
	},
	focusToEnd : function() {
        return this.each(function() {
            var v = $(this).val();
            $(this).focus().val("").val(v);
        });
    }
});

var chatMsgTemplate = function(msg, name, commentType, sinceTime, id, parentC) {
	var replyBtnHtml = '';
	var editBtn= '';
	if(parentC == "parent") {
		replyBtnHtml = ' <i class="fa fa-reply replyBtn" aria-hidden="true"></i>'
	}
	if(commentType) {
		editBtn = '<i class="fa fa-pencil editCommBtn" aria-hidden="true"></i>';
	}
	var html = '<li class="media '+commentType+' '+parentC+'" data-id="'+id+'">'+
		'<div class="media-body">'+
			'<span class="media-annotation display-block"><b>'+name+'</b>, '+sinceTime+replyBtnHtml+" "+editBtn+'</span>'+
			'<div class="media-content">'+msg+'</div>'+
		'</div>'+
	'</li>';
	return html;
}

function commentReply(_this) {
	var postID = _this.closest("li").data("id");
	$("#atjs").html("@"+postID+"&nbsp;");//.focusToEnd();
	var el = document.getElementById('atjs')
	focusAndPlaceCaretAtEnd(el);
}

function timeSince(date) {
  if (typeof date !== 'object') {
    date = new Date(date);
  }

  var seconds = Math.floor((new Date() - date) / 1000);
  var intervalType;

  var interval = Math.floor(seconds / 31536000);
  if (interval >= 1) {
    intervalType = 'year';
  } else {
    interval = Math.floor(seconds / 2592000);
    if (interval >= 1) {
      intervalType = 'month';
    } else {
      interval = Math.floor(seconds / 86400);
      if (interval >= 1) {
        intervalType = 'day';
      } else {
        interval = Math.floor(seconds / 3600);
        if (interval >= 1) {
          intervalType = "hr";
        } else {
          interval = Math.floor(seconds / 60);
          if (interval >= 1) {
            intervalType = "min";
          } else {
            interval = seconds;
            intervalType = "sec";
          }
        }
      }
    }
  }

  if (interval > 1 || interval <= 0 ) {
	if(intervalType == "sec") {
		return "Just now";
	}
	else if (intervalType == "day" || intervalType == "month" || intervalType == "year") {
		return date.toLocaleString();
	}
    intervalType += 's';
  }

  return interval + ' ' + intervalType;
};

function getCommentID(msg) {
	var parentID = msg.trim().split(' ')[0].replace(/\s/g, '');
	var first = parentID.charAt(0);
	if(first === "@" && parentID.length == "25") {
		return parentID.slice(1);
	}
}

function editComm(_this) {
	editCommentFlag = true;
	$(".media-list .media.edit").removeClass("edit");
	
	_this.closest("li").addClass("edit");
	var postText = _this.closest("li").find(".media-content").html();
	$("#atjs").html("~"+postText);//.focusToEnd();
	var el = document.getElementById('atjs')
	focusAndPlaceCaretAtEnd(el);
}


function isEditActive(msg) {
	var first = msg.charAt(0);
	if(editCommentFlag && first=="~") {
		return true;
	}
	return false;
}

function updateComm(msg) {
	//var msg = msg.slice(1);
	var postID = $(".media.edit").data("id");
	var data = {};
	data.postID = postID;
	data.msg = msg;
	$("#atjs").html("");
	$.ajax({
		type:"PUT",
		url:"/business/comment",
		data: data,
		success: function(data) {
			if(data.success) {
				$("li[data-id='"+postID+"'] .media-content").html(msg);
				socket.OnConnect(function () {
					socket.Emit("refreshCommBiz", bizID);
				});
			}
			$(".media-list .media.edit").removeClass("edit");
			editCommentFlag = false;
		},
		dataType: 'json',
	});
}

function settingsTable() {
	var dt2 = $activityTypeTable.DataTable( {
        "processing": true,
		"serverSide": true,
		"keys": true,
		"autoWidth": false,
		"stateSave": false,
		"order": [[1, "asc"]],
        "ajax": "/activitylist",
		"rowId": "id",
		"columns": [
			{
				"data": null,
                "defaultContent": "",
				"width": "50",
				"orderable": false,
			},
			{ 
				name: "activityName",
				data: "name",
				render : function(data, type, row) {
					return '<span class="editActivity">'+data+'</span>'
				}    
			},
			{
				"data": null,
				//"class": '',
                "defaultContent": '<i class="icon-trash deleteActivityType"></i>',
				"width": '50',
				"orderable": false,
			}
		],
		"dom": '<"datatable-header"fl><"datatable-scroll-wrap"t><"datatable-footer"ip>',
        "language": {
            search: '<span>Filter:</span> _INPUT_',
            searchPlaceholder: 'Type to filter...',
            lengthMenu: '<span>Show:</span> _MENU_',
            paginate: { 'first': 'First', 'last': 'Last', 'next': '&rarr;', 'previous': '&larr;' }
        },
		"drawCallback": function () {
			activityTypeEdit();
			dt2.column(0, {search:'applied', order:'applied'}).nodes().each( function (cell, i) {
				var start = dt2.page.info().start;
				cell.innerHTML = start+i+1;
			});
		}
    });
	



	$('.activityTypeNewBtn').on('click', function() {
		addNewActivity(dt2);
	});
	
	$(".activityType").on('keydown', function(event) {
		if (event.keyCode == 13) {
			addNewActivity(dt2);
		}
	})
	
	
	$('#activityTypeTable').on('click', '.deleteActivityType', function() {
		var activityTypeID = $(this).closest('tr').attr("id");
		$.ajax({
			type:"DELETE",
			url:"/activitytype/"+activityTypeID,
			success: function(data) {
				if(data.success) {
					dt2.ajax.reload(null, false);
					new PNotify({
						title: 'Value Deleted',
						text: 'Value was deteled successfully!',
						type: 'success',
						icon: 'icon-checkmark3'
					});
				} else {
					new PNotify({
						title: 'Value was not deleted!',
						text: 'There was an error deleting the value, please try again.',
						type: 'error',
						icon: 'icon-blocked'
					});
				}
				
			},
			dataType: 'json',
		});	
	});
}

function addNewActivity(dt2) {
	var name = $(".activityType").val();
	$(".activityType").focus();
	if(name != "") {
		$(".activityType").val("");
		$.ajax({
			type:"POST",
			url:"/activitytype",
			data: {"name": name},
			success: function(data) {
				dt2.ajax.reload(null, false);
				new PNotify({
					title: 'Value Added successfully',
					text: 'Value was added successfully!',
					type: 'success',
					icon: 'icon-checkmark3'
				});
			},
			dataType: 'json',
		});	
	}
}

function activityTypeEdit() {
	$('#activityTypeTable .editActivity').editable({
		url: '/activitylist',
		type: 'text',
		params: function(params) {
			params.pk = $(this).closest('tr').attr("id");
			return params;
		},
		validate: function(value) {
            if($.trim(value) == '') return 'This field is required';
        },
		ajaxOptions:{
			type:'put',
			dataType: 'json'
		} ,
		pk: '',
		name: 'name',
		mode: 'inline',
		success: function(response, newValue) {
			if(response.success) {
				$(this).html(newValue);
				new PNotify({
					title: 'Value Updated',
					text: 'Value was updatedd successfully!',
					type: 'success',
					icon: 'icon-checkmark3'
				});
			} else {
			new PNotify({
				title: 'Value was not updated!',
				text: 'There was an error, please try again.',
				type: 'error',
				icon: 'icon-blocked'
			});
			}
		}
	});
}

function focusAndPlaceCaretAtEnd(el) {
    el.focus();
    if (typeof window.getSelection != "undefined"
            && typeof document.createRange != "undefined") {
        var range = document.createRange();
        range.selectNodeContents(el);
        range.collapse(false);
        var sel = window.getSelection();
        sel.removeAllRanges();
        sel.addRange(range);
        
    } else if (typeof document.body.createTextRange != "undefined") {
        var textRange = document.body.createTextRange();
        textRange.moveToElementText(el);
        textRange.collapse(false);
        textRange.select();
    }
}


function highChartsAll() {
	
	var option = {

        chart: {
            renderTo: 'flipbookStats',
            zoomType: 'x',
			type: 'column'
        },
		
		 title: {
            text: 'Statistics per week'
        },
		 subtitle: {
            text: 'Click and drag to zoom in' 
        },
		 credits: {
			  enabled: false
		 },
		tooltip: {
            shared: true,
            crosshairs: true
        },
		yAxis: {
            title: {
                text: 'Activity Count'
            }
        },
        xAxis: {
            type: 'category',
        },
		series: []
    
    }
	
	var data = dataHighchartsActivityType(activityWeek);

	option.series = data;
	var chart = new Highcharts.Chart(option);
	
	$('.selectAdmin').on('input select2:select', function() {
		var _this = $(this);
		var adminID = _this.val();
		if(adminID != "") {
			$.ajax({
				type:"GET",
				url:"/activitytype/"+adminID,
				success: function(data) {
					console.log(chart);
					if(!isEmpty(chart)) {
						chart.destroy();
					}
					var stats = dataHighchartsActivityType(data.data);
					if(!stats.length) {
						$("#flipbookStats").html("<h3 class='text-center'>No statistics available for this user yet!</h3>");
						return;
					}
					option.series = stats;
					chart = new Highcharts.Chart(option);
				},
				dataType: 'json',
			});
		}
	});
}

function dataHighchartsActivityType(dataRaw) {
	var data = [];
	var nr = 0;
	var found = true;
	for(var i = 0; i < dataRaw.length; i++) {
		var item = dataRaw[i];
		
		for(var j= 0; j < item.activity_types.length; j++) {
			var subitem = item.activity_types[j];
			if(data.length) {
				for(var k= 0; k < data.length; k++) {
					var subSubItem = data[k];
					if(subitem.name == subSubItem.name) {
						found = true;
						subSubItem.data.push(["week "+item._id.week+","+item._id.year, subitem.sum]);
						break;
					} else {	
						found = false;

					}
				}
			} else {
				data[nr] = {
					name: subitem.name,
					data: []
				};
				data[nr].data.push(["week "+item._id.week+","+item._id.year, subitem.sum]);
				nr++;
			}
			if(!found) {
				data[nr] = {
					name: subitem.name,
					data: []
				};
				data[nr].data.push(["week "+item._id.week+","+item._id.year, subitem.sum]);
				nr++;
			}
		}
		/*if(item._id === null) {
			item._id = "undefined";
		}
		data[i] = {
			name: item._id,
			data: []
		};
		for(var j= 0; j < item.activity_types.length; j++) {
			data[i].data.push([item.activity_types[j].name+","+item.activity_types[j].name, item.activity_types[j].sum]);
			
		}*/
		
	}
	return data;
}

  var justifyColumns = function(chart) {
    var categoriesWidth = chart.plotSizeX / (1 + chart.xAxis[0].max - chart.xAxis[0].min),
      distanceBetweenColumns = 0,
      each = Highcharts.each,
      sum, categories = chart.xAxis[0].categories,
      number;
    for (var i = 0; i < categories.length; i++) {
      sum = 0;
      each(chart.series, function(p, k) {
        if (p.visible) {
          each(p.data, function(ob, j) {
            if (ob.category == categories[i]) {
              sum++;
            }
          });
        }
      });
      distanceBetweenColumns = categoriesWidth / (sum + 1);
      number = 1;
      each(chart.series, function(p, k) {
        if (p.visible) {
          each(p.data, function(ob, j) {
            if (ob.category == categories[i] && typeof(ob.graphic) !== 'undefined') {
              ob.graphic.element.x.baseVal.value = i * categoriesWidth + distanceBetweenColumns * number - ob.pointWidth / 2;
              number++;
            }
          });
        }
      });
    }
  };
  
function isEmpty(obj) {
    for(var key in obj) {
        if(obj.hasOwnProperty(key))
            return false;
    }
    return true;
}