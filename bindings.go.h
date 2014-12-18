#include <stdlib.h>
#include <glib.h>
#include <glib-object.h>
#include <gtk/gtk.h>

static void
_g_object_set_one(gpointer object, const gchar *property_name, void *val)
{
	g_object_set(object, property_name, *(gpointer **)val, NULL);
}
