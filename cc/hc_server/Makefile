default: hc_server

CXX=clang++
LD=clang++

DEFINES:=-DV8_DEPRECATION_WARNINGS -DENABLE_MDNS=1 -DENABLE_NOTIFICATIONS -DENABLE_PEPPER_CDMS \
	-DENABLE_PLUGINS=1 -DENABLE_PDF=1 -DENABLE_SPELLCHECK=1 -DUI_COMPOSITOR_IMAGE_TRANSPORT \
	-DUSE_AURA=1 -DUSE_DEFAULT_RENDER_THEME=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DENABLE_WEBRTC=1 \
	-DDISABLE_NACL -DENABLE_EXTENSIONS=1 -DENABLE_TASK_MANAGER=1 -DENABLE_THEMES=1 \
	-DENABLE_CAPTIVE_PORTAL_DETECTION=1 -DENABLE_SESSION_SERVICE=1 -DENABLE_SUPERVISED_USERS=1 \
	-DFULL_SAFE_BROWSING -DSAFE_BROWSING_CSD -DSAFE_BROWSING_DB_LOCAL -DCHROMIUM_BUILD \
	-DENABLE_MEDIA_ROUTER=1 -DFIELDTRIAL_TESTING_ENABLED -D_FILE_OFFSET_BITS=64 \
	-D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS \
	-DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DHEADLESS_IMPLEMENTATION -DUSE_EGL \
	-DENABLE_IPC_FUZZER -DSK_IGNORE_DW_GRAY_FIX -DSK_IGNORE_LINEONLY_AA_CONVEX_PATH_OPTS \
	-DSK_SUPPORT_GPU=1 -DU_USING_ICU_NAMESPACE=0 -DU_ENABLE_DYLOAD=0 -DU_NOEXCEPT= \
	-DU_STATIC_IMPLEMENTATION -DICU_UTIL_DATA_IMPL=ICU_UTIL_DATA_STATIC -DENABLE_WEBSOCKETS \
	-DGOOGLE_PROTOBUF_NO_RTTI -DGOOGLE_PROTOBUF_NO_STATIC_INITIALIZER -DHAVE_PTHREAD \
	-DMESA_EGL_NO_X11_HEADERS

CXXFLAGS:=-fno-strict-aliasing --param=ssp-buffer-size=4 -fstack-protector -funwind-tables -fPIC \
	-pipe -pthread -m64 -march=x86-64 -O2 -fno-ident -fdata-sections -ffunction-sections -g0 \
	-fno-strict-overflow -fvisibility=hidden  -fno-threadsafe-statics -std=gnu++11 -fno-rtti \
	-fno-exceptions -I/usr/local/headless_chromium/include

ifdef DEBUG_CHROMIUM
	CHROMIUM_LIBS:=-laccessibility -laura -laura_extra -lbase -lbase_i18n -lbitmap_uploader \
		-lblink_common -lblink_core -lblink_modules -lblink_platform -lblink_web -lbluetooth \
		-lboringssl -lcc -lcc_blink -lcc_ipc -lcc_proto -lcc_surfaces -lchromium_sqlite3 \
		-lclient_sources_for_ppapi -lcompositor -lcontent -lcrcrypto -ldevice_battery \
		-ldevice_vibration -ldevices -ldisplay -ldisplay_compositor -ldisplay_types \
		-ldisplay_util -levents -levents_base -levents_ipc -levents_ozone -levents_ozone_evdev \
		-levents_ozone_layout -lffmpeg -lgeometry -lgesture_detection -lgfx -lgfx_ipc \
		-lgfx_ipc_geometry -lgfx_ipc_skia -lgin -lgl_init -lgl_wrapper -lgles2_c_lib \
		-lgles2_implementation -lgles2_utils -lgpu -licui18n -licuuc -lipc -lmedia -lmedia_blink \
		-lmedia_gpu -lmidi -lmojo -lmojo_blink_lib -lmojo_common_lib -lmojo_ime_lib \
		-lmojo_public_system -lmojo_surfaces_lib -lmojo_system_impl -lmus_common -lmus_library \
		-lnative_theme -lnet -lozone -lozone_base -lplatform -lppapi_host -lppapi_proxy \
		-lppapi_shared -lprefs -lprotobuf_lite -lrange -lsandbox_services -lscheduler \
		-lseccomp_bpf -lseccomp_bpf_helpers -lshared_memory_support -lshell_dialogs -lskia \
		-lsnapshot -lsql -lstartup_tracing -lstorage_browser -lstorage_common -lstub_window \
		-lsuid_sandbox_client -lsurface -ltracing -ltracing_library -ltranslator -lui_base \
		-lui_base_ime -lui_data_pack -lui_touch_selection -lurl -lurl_ipc -lv8 -lwm -lwtf
else
	CHROMIUM_LIBS:=-lheadless_lib
endif

LDFLAGS:=-L/usr/local/headless_chromium/lib $(CHROMIUM_LIBS) -lnss3 -lnssutil3 -lsmime3 -lplds4 \
	-lplc4 -lnspr4 -latomic -lfontconfig -lfreetype -lexpat -lresolv -lm -lz -ldl -lrt -lpthread

browser.o: browser.cc browser.hpp
	$(CXX) -c $(DEFINES) $(CXXFLAGS) $< -o $@

main.o: main.cc browser.hpp
	$(CXX) -c $(DEFINES) $(CXXFLAGS) $< -o $@

hc_server: browser.o main.o
	$(LD) -o $@ $^ $(LDFLAGS)

.PHONY: clean

clean:
	rm -f *.o hc_server
