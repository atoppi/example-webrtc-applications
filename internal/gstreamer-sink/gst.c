// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

#include "gst.h"

#include <gst/app/gstappsrc.h>

GMainLoop *gstreamer_receive_main_loop = NULL;
void gstreamer_receive_start_mainloop(void) {
  gstreamer_receive_main_loop = g_main_loop_new(NULL, FALSE);

  g_main_loop_run(gstreamer_receive_main_loop);
}

static gboolean gstreamer_receive_bus_call(GstBus *bus, GstMessage *msg, gpointer data) {
  switch (GST_MESSAGE_TYPE(msg)) {

  case GST_MESSAGE_EOS:
    g_print("End of stream\n");
    exit(1);
    break;

  case GST_MESSAGE_ERROR: {
    GError *err = NULL;
    gchar *dbg_info = NULL;

    gst_message_parse_error (msg, &err, &dbg_info);
    g_printerr("ERROR from element %s: %s\n", GST_OBJECT_NAME (msg->src), err->message);
    g_printerr("Debugging info: %s\n", (dbg_info) ? dbg_info : "none");
    g_error_free(err);
    g_free(dbg_info);
    exit(1);
  }
  default:
    break;
  }

  return TRUE;
}

GstElement *gstreamer_receive_create_pipeline(char *pipeline) {
  gst_init(NULL, NULL);
  GError *error = NULL;
  return gst_parse_launch(pipeline, &error);
}

void gstreamer_receive_start_pipeline(GstElement *pipeline) {
  GstBus *bus = gst_pipeline_get_bus(GST_PIPELINE(pipeline));
  gst_bus_add_watch(bus, gstreamer_receive_bus_call, NULL);
  gst_object_unref(bus);

  gst_element_set_state(pipeline, GST_STATE_PLAYING);
}

void gstreamer_receive_stop_pipeline(GstElement *pipeline) { gst_element_set_state(pipeline, GST_STATE_NULL); }

void gstreamer_receive_push_buffer(GstElement *pipeline, void *buffer, int len) {
  GstElement *src = gst_bin_get_by_name(GST_BIN(pipeline), "src");
  if (src != NULL) {
    gpointer p = g_memdup2(buffer, len);
    GstBuffer *buffer = gst_buffer_new_wrapped(p, len);
    gst_app_src_push_buffer(GST_APP_SRC(src), buffer);
    gst_object_unref(src);
  }
}
