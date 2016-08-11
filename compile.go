package headless_chromium

/*
#cgo CXXFLAGS: -DCHROMIUM_BUILD -DDISABLE_NACL -DENABLE_CAPTIVE_PORTAL_DETECTION=1
#cgo CXXFLAGS: -DENABLE_EXTENSIONS=1 -DENABLE_MEDIA_ROUTER=1 -DENABLE_MDNS=1
#cgo CXXFLAGS: -DENABLE_NOTIFICATIONS -DENABLE_PDF=1 -DENABLE_PEPPER_CDMS -DENABLE_PLUGINS=1
#cgo CXXFLAGS: -DENABLE_SERVICE_DISCOVERY=1 -DENABLE_SESSION_SERVICE=1 -DENABLE_SPELLCHECK=1
#cgo CXXFLAGS: -DENABLE_SUPERVISED_USERS=1 -DENABLE_TASK_MANAGER=1 -DENABLE_THEMES=1
#cgo CXXFLAGS: -DENABLE_WEBRTC=1 -DFULL_SAFE_BROWSING -DSAFE_BROWSING_CSD -DSAFE_BROWSING_DB_LOCAL
#cgo CXXFLAGS: -DUI_COMPOSITOR_IMAGE_TRANSPORT -DUSE_AURA=1 -DUSE_DEFAULT_RENDER_THEME=1
#cgo CXXFLAGS: -DUSE_EGL -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DV8_DEPRECATION_WARNINGS

#cgo CXXFLAGS: -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FILE_OFFSET_BITS=64
#cgo CXXFLAGS: -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -DDYNAMIC_ANNOTATIONS_ENABLED=1
#cgo CXXFLAGS: -DFIELDTRIAL_TESTING_ENABLED -DMOJO_USE_SYSTEM_IMPL -DTOOLKIT_VIEWS=1
#cgo CXXFLAGS: -DWTF_USE_DYNAMIC_ANNOTATIONS=1

#cgo CXXFLAGS: -fno-exceptions -fno-rtti -fno-strict-aliasing -fno-threadsafe-statics
#cgo CXXFLAGS: -fstack-protector -funwind-tables -fvisibility-inlines-hidden

#cgo CXXFLAGS: -g -O2 -DNDEBUG -Wall
#cgo CXXFLAGS: -std=gnu++11 -gsplit-dwarf --param=ssp-buffer-size=4 -pipe -pthread

#cgo CXXFLAGS: -I/usr/local/include/headless_chromium

#cgo LDFLAGS: -lheadless_chromium -lnss3 -lnssutil3 -lsmime3 -lplds4 -lplc4 -lnspr4
#cgo LDFLAGS: -lfontconfig -lfreetype -lexpat -lresolv -lm -lz -ldl -lrt
#cgo LDFLAGS: -laccessibility -laura -laura_extra -lbase -lbase_i18n -lbitmap_uploader -lblink_common -lblink_core -lblink_modules -lblink_platform -lblink_web -lbluetooth -lboringssl -lcc -lcc_blink -lcc_ipc -lcc_proto -lcc_surfaces -lchromium_sqlite3 -lclient_sources_for_ppapi -lcompositor -lcontent -lcrcrypto -ldevice_battery -ldevice_vibration -ldevices -ldisplay -ldisplay_compositor -ldisplay_types -ldisplay_util -levents -levents_base -levents_ipc -levents_ozone -levents_ozone_evdev -levents_ozone_layout -lffmpeg -lgeometry -lgesture_detection -lgfx -lgfx_ipc -lgfx_ipc_geometry -lgfx_ipc_skia -lgin -lgl_init -lgl_wrapper -lgles2_c_lib -lgles2_implementation -lgles2_utils -lgpu -licui18n -licuuc -lipc -lmedia -lmedia_blink -lmedia_gpu -lmidi -lmojo -lmojo_blink_lib -lmojo_common_lib -lmojo_ime_lib -lmojo_public_system -lmojo_surfaces_lib -lmojo_system_impl -lmus_common -lmus_library -lnative_theme -lnet -lozone -lozone_base -lplatform -lppapi_host -lppapi_proxy -lppapi_shared -lprefs -lprotobuf_lite -lrange -lsandbox_services -lscheduler -lseccomp_bpf -lseccomp_bpf_helpers -lshared_memory_support -lshell_dialogs -lskia -lsnapshot -lsql -lstartup_tracing -lstorage_browser -lstorage_common -lstub_window -lsuid_sandbox_client -lsurface -ltracing -ltracing_library -ltranslator -lui_base -lui_base_ime -lui_data_pack -lui_touch_selection -lurl -lurl_ipc -lv8 -lwm -lwtf
*/
import "C"
