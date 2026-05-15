package main

/*
#cgo pkg-config: libgpod-1.0 glib-2.0
#include <gpod/itdb.h>
#include <glib.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

typedef struct {
	char   *model;
	char   *guid;
	char   *modelnum;
	guint64 capacity;
	int     track_count;
	int     connected;
} DeviceInfo;

typedef struct {
	char *title;
	char *artist;
	char *album;
	int   id;
} TrackInfo;

static int has_storage_info() {
#ifdef ITDB_DEVICE_GET_STORAGE_INFO
	return 1;
#else
	return 0;
#endif
}

DeviceInfo get_device_info(const char *mountpoint) {
	DeviceInfo info = {0};
	GError *err = NULL;
	Itdb_iTunesDB *itdb = itdb_parse(mountpoint, &err);
	if (!itdb) {
		if (err) g_error_free(err);
		return info;
	}
	info.connected   = 1;
	info.track_count = itdb_tracks_number(itdb);
	Itdb_Device *dev = itdb->device;
	if (dev) {
		const Itdb_IpodInfo *ipod_info = itdb_device_get_ipod_info(dev);
		if (ipod_info) {
			info.model    = g_strdup(itdb_info_get_ipod_model_name_string(ipod_info->ipod_model));
			info.capacity = (guint64)(ipod_info->capacity * 1024 * 1024 * 1024);
		}
		info.guid     = g_strdup(itdb_device_get_sysinfo(dev, "FirewireGuid"));
		info.modelnum = g_strdup(itdb_device_get_sysinfo(dev, "ModelNumStr"));
#ifdef ITDB_DEVICE_GET_STORAGE_INFO
		double cap_bytes = 0, free_bytes = 0;
		itdb_device_get_storage_info(dev, &cap_bytes, &free_bytes);
		if (cap_bytes > 0)
			info.capacity = (guint64)cap_bytes;
#endif
	}
	itdb_free(itdb);
	return info;
}

void free_device_info(DeviceInfo *info) {
	if (info->model)    g_free(info->model);
	if (info->guid)     g_free(info->guid);
	if (info->modelnum) g_free(info->modelnum);
}

TrackInfo *get_tracks_info(const char *mountpoint, int *out_count) {
	GError *err = NULL;
	Itdb_iTunesDB *itdb = itdb_parse(mountpoint, &err);
	if (!itdb) {
		if (err) g_error_free(err);
		*out_count = 0;
		return NULL;
	}
	int count = itdb_tracks_number(itdb);
	TrackInfo *infos = (TrackInfo *)malloc(sizeof(TrackInfo) * count);
	int i = 0;
	GList *tl = itdb->tracks;
	while (tl) {
		Itdb_Track *t = (Itdb_Track *)tl->data;
		infos[i].title  = g_strdup((t->title  && strlen(t->title))  ? t->title  : "Unknown Title");
		infos[i].artist = g_strdup((t->artist && strlen(t->artist)) ? t->artist : "Unknown Artist");
		infos[i].album  = g_strdup((t->album  && strlen(t->album))  ? t->album  : "Unknown Album");
		infos[i].id     = t->id;
		i++;
		tl = tl->next;
	}
	*out_count = count;
	itdb_free(itdb);
	return infos;
}

void free_tracks_info(TrackInfo *infos, int count) {
	if (!infos) return;
	for (int i = 0; i < count; i++) {
		g_free(infos[i].title);
		g_free(infos[i].artist);
		g_free(infos[i].album);
	}
	free(infos);
}

int add_track(const char *mountpoint, const char *file_path) {
	GError *err = NULL;
	Itdb_iTunesDB *itdb = itdb_parse(mountpoint, &err);
	if (!itdb) {
		if (err) { g_printerr("parse failed: %s\n", err->message); g_error_free(err); }
		return -1;
	}
	Itdb_Track *track = itdb_track_new();
	track->mediatype = 0x00000001;
	char *base = g_path_get_basename(file_path);
	char *dot  = strrchr(base, '.');
	if (dot) *dot = '\0';
	track->title  = g_strdup(base);
	track->artist = g_strdup("Unknown Artist");
	track->album  = g_strdup("Unknown Album");
	g_free(base);
	itdb_track_add(itdb, track, -1);
	Itdb_Playlist *mpl = itdb_playlist_mpl(itdb);
	if (mpl) itdb_playlist_add_track(mpl, track, -1);
	if (!itdb_cp_track_to_ipod(track, (char *)file_path, &err)) {
		if (err) { g_printerr("cp failed: %s\n", err->message); g_error_free(err); }
		itdb_free(itdb);
		return -2;
	}
	if (!itdb_write(itdb, &err)) {
		if (err) { g_printerr("write failed: %s\n", err->message); g_error_free(err); }
		itdb_free(itdb);
		return -3;
	}
	itdb_free(itdb);
	return 0;
}

int delete_track_by_id(const char *mountpoint, int track_id) {
	GError *err = NULL;
	Itdb_iTunesDB *itdb = itdb_parse(mountpoint, &err);
	if (!itdb) {
		if (err) g_error_free(err);
		return -1;
	}
	Itdb_Track *track = itdb_track_by_id(itdb, (guint32)track_id);
	if (!track) {
		itdb_free(itdb);
		return -2;
	}
	char *ipod_path = itdb_filename_on_ipod(track);
	GList *pl = itdb->playlists;
	while (pl) {
		Itdb_Playlist *p = (Itdb_Playlist *)pl->data;
		if (itdb_playlist_contains_track(p, track))
			itdb_playlist_remove_track(p, track);
		pl = pl->next;
	}
	itdb_track_remove(track);
	if (!itdb_write(itdb, &err)) {
		if (err) { g_printerr("write failed: %s\n", err->message); g_error_free(err); }
		g_free(ipod_path);
		itdb_free(itdb);
		return -3;
	}
	if (ipod_path) {
		if (remove(ipod_path) != 0)
			g_printerr("warning: could not delete file %s\n", ipod_path);
		g_free(ipod_path);
	}
	itdb_free(itdb);
	return 0;
}
*/
import "C"

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TrackInfo struct {
	Title  string
	Artist string
	Album  string
	ID     int
}

type DeviceInfo struct {
	Connected  bool
	Model      string
	GUID       string
	ModelNum   string
	CapacityGB float64
	UsedGB     float64
	FreeGB     float64
	TrackCount int
}

var (
	mu                sync.Mutex
	currentMountpoint string
	fullTracks        []TrackInfo
	filteredIndices   []int
	selected          = make(map[int]bool)
)

func detectMountpoint() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	base := filepath.Join("/run/media", u.Username)
	entries, err := os.ReadDir(base)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		candidate := filepath.Join(base, e.Name())
		if _, err := os.Stat(filepath.Join(candidate, "iPod_Control")); err == nil {
			return candidate
		}
	}
	return ""
}

func diskUsage(mp string) (capGB, usedGB, freeGB float64) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(mp, &stat); err != nil {
		return 0, 0, 0
	}
	gb    := 1024.0 * 1024.0 * 1024.0
	cap_  := float64(stat.Blocks) * float64(stat.Bsize)
	avail := float64(stat.Bavail) * float64(stat.Bsize)
	return cap_ / gb, (cap_ - avail) / gb, avail / gb
}

func isAudio(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp3", ".aac", ".m4a", ".flac", ".ogg", ".wav":
		return true
	}
	return false
}

func fetchDeviceInfo(mp string) DeviceInfo {
	cMp := C.CString(mp)
	defer C.free(unsafe.Pointer(cMp))
	raw := C.get_device_info(cMp)
	defer C.free_device_info(&raw)
	if raw.connected == 0 {
		return DeviceInfo{}
	}
	strOr := func(p *C.char, def string) string {
		if p != nil {
			return C.GoString(p)
		}
		return def
	}
	capGB, usedGB, freeGB := diskUsage(mp)
	if C.has_storage_info() == 1 {
		capGB = float64(raw.capacity) / (1024 * 1024 * 1024)
	}
	return DeviceInfo{
		Connected:  true,
		Model:      strOr(raw.model, "Unknown"),
		GUID:       strOr(raw.guid, "-"),
		ModelNum:   strOr(raw.modelnum, "-"),
		CapacityGB: capGB,
		UsedGB:     usedGB,
		FreeGB:     freeGB,
		TrackCount: int(raw.track_count),
	}
}

func fetchTracksInfo(mp string) []TrackInfo {
	cMp := C.CString(mp)
	defer C.free(unsafe.Pointer(cMp))
	var count C.int
	raw := C.get_tracks_info(cMp, &count)
	if raw == nil {
		return nil
	}
	defer C.free_tracks_info(raw, count)
	type cTrack struct {
		title  *C.char
		artist *C.char
		album  *C.char
		id     C.int
	}
	slice := unsafe.Slice((*cTrack)(unsafe.Pointer(raw)), int(count))
	tracks := make([]TrackInfo, int(count))
	for i, t := range slice {
		tracks[i] = TrackInfo{
			Title:  C.GoString(t.title),
			Artist: C.GoString(t.artist),
			Album:  C.GoString(t.album),
			ID:     int(t.id),
		}
	}
	return tracks
}

func doAddTrack(mp, path string) error {
	cMp   := C.CString(mp)
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cMp))
	defer C.free(unsafe.Pointer(cPath))
	if ret := C.add_track(cMp, cPath); ret != 0 {
		return fmt.Errorf("code %d", int(ret))
	}
	return nil
}

func doDeleteTrack(mp string, id int) error {
	cMp := C.CString(mp)
	defer C.free(unsafe.Pointer(cMp))
	if ret := C.delete_track_by_id(cMp, C.int(id)); ret != 0 {
		return fmt.Errorf("code %d", int(ret))
	}
	return nil
}

func main() {
	a := app.NewWithID("com.fcnst.podthing")
	w := a.NewWindow("PodThing - Alpha")
	w.Resize(fyne.NewSize(800, 500))
	w.CenterOnScreen()

	statusLabel  := widget.NewLabel("○ Disconnected")
	statusLabel.TextStyle = fyne.TextStyle{Bold: true}
	modelVal     := widget.NewLabel("-")
	guidVal      := widget.NewLabel("-")
	capacityVal  := widget.NewLabel("-")
	usedVal      := widget.NewLabel("-")
	freeVal      := widget.NewLabel("-")
	songsVal     := widget.NewLabel("-")
	progressLabel := widget.NewLabel("")

	updateDeviceUI := func(info DeviceInfo) {
		if !info.Connected {
			statusLabel.SetText("○ Disconnected")
			modelVal.SetText("-")
			guidVal.SetText("-")
			capacityVal.SetText("-")
			usedVal.SetText("-")
			freeVal.SetText("-")
			songsVal.SetText("-")
			return
		}
		statusLabel.SetText("● Connected")
		modelVal.SetText(info.Model)
		guidVal.SetText(fmt.Sprintf("%s", info.GUID))
		capacityVal.SetText(fmt.Sprintf("%.1f GB", info.CapacityGB))
		usedVal.SetText(fmt.Sprintf("%.1f GB", info.UsedGB))
		freeVal.SetText(fmt.Sprintf("%.1f GB", info.FreeGB))
		songsVal.SetText(fmt.Sprintf("%d", info.TrackCount))
	}

	refreshStorage := func(mp string) {
		if mp == "" {
			return
		}
		_, usedGB, freeGB := diskUsage(mp)
		info := fetchDeviceInfo(mp)
		fyne.Do(func() {
			usedVal.SetText(fmt.Sprintf("%.1f GB", usedGB))
			freeVal.SetText(fmt.Sprintf("%.1f GB", freeGB))
			songsVal.SetText(fmt.Sprintf("%d", info.TrackCount))
		})
	}

	var songList *widget.List
	songList = widget.NewList(
		func() int { return len(filteredIndices) },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewCheck("", nil), widget.NewLabel(""))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box   := obj.(*fyne.Container)
			check := box.Objects[0].(*widget.Check)
			label := box.Objects[1].(*widget.Label)
			if id >= len(filteredIndices) {
				return
			}
			t := fullTracks[filteredIndices[id]]
			label.SetText(t.Title)
			check.SetChecked(selected[id])
			check.OnChanged = func(v bool) { selected[id] = v }
		},
	)

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search by title, artist or album...")

	rebuildList := func() {
		filter := strings.ToLower(searchEntry.Text)
		filteredIndices = filteredIndices[:0]
		for i, t := range fullTracks {
			if filter == "" ||
				strings.Contains(strings.ToLower(t.Title), filter) ||
				strings.Contains(strings.ToLower(t.Artist), filter) ||
				strings.Contains(strings.ToLower(t.Album), filter) {
				filteredIndices = append(filteredIndices, i)
			}
		}
		for k := range selected {
			delete(selected, k)
		}
		songList.Refresh()
	}

	searchEntry.OnChanged = func(string) { rebuildList() }

	refreshTracks := func(mp string) {
		var tracks []TrackInfo
		if mp != "" {
			tracks = fetchTracksInfo(mp)
		}
		fyne.Do(func() {
			mu.Lock()
			fullTracks = tracks
			mu.Unlock()
			rebuildList()
		})
	}

	toggleBtn    := widget.NewButton("Connect", nil)
	addBtn       := widget.NewButtonWithIcon("Add Songs", theme.FolderOpenIcon(), nil)
	addFolderBtn := widget.NewButtonWithIcon("Add Folder", theme.FolderIcon(), nil)
	deleteBtn    := widget.NewButtonWithIcon("Delete Selected", theme.DeleteIcon(), nil)
	deleteBtn.Importance = widget.DangerImportance

	setBusy := func(b bool) {
		if b {
			addBtn.Disable()
			addFolderBtn.Disable()
			deleteBtn.Disable()
			toggleBtn.Disable()
		} else {
			addBtn.Enable()
			addFolderBtn.Enable()
			deleteBtn.Enable()
			toggleBtn.Enable()
			progressLabel.SetText("")
		}
	}

	runBatch := func(mp string, paths []string, op string, fn func(string) error) {
		total := len(paths)
		var errs []string
		for i, p := range paths {
			idx := i + 1
			fyne.Do(func() {
				progressLabel.SetText(fmt.Sprintf("%s %d / %d", op, idx, total))
			})
			if err := fn(p); err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", filepath.Base(p), err))
			}
		}
		go refreshStorage(mp)
		refreshTracks(mp)
		fyne.Do(func() {
			setBusy(false)
			if len(errs) > 0 {
				dialog.ShowError(fmt.Errorf("%s", strings.Join(errs, "\n")), w)
			}
		})
	}

	toggleBtn.OnTapped = func() {
		if toggleBtn.Text == "Connect" {
			mp := detectMountpoint()
			if mp == "" {
				dialog.ShowError(fmt.Errorf("no iPod found under /run/media"), w)
				return
			}
			mu.Lock()
			currentMountpoint = mp
			mu.Unlock()
			info := fetchDeviceInfo(mp)
			updateDeviceUI(info)
			if info.Connected {
				toggleBtn.SetText("Disconnect")
				setBusy(true)
				go func() {
					refreshTracks(mp)
					fyne.Do(func() { setBusy(false) })
				}()
			}
		} else {
			mu.Lock()
			currentMountpoint = ""
			mu.Unlock()
			updateDeviceUI(DeviceInfo{})
			go refreshTracks("")
			toggleBtn.SetText("Connect")
		}
	}

	addBtn.OnTapped = func() {
		mu.Lock()
		mp := currentMountpoint
		mu.Unlock()
		if mp == "" {
			dialog.ShowError(fmt.Errorf("no iPod connected"), w)
			return
		}
		fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			defer uc.Close()
			path := uc.URI().Path()
			setBusy(true)
			go runBatch(mp, []string{path}, "Adding", func(p string) error {
				return doAddTrack(mp, p)
			})
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{
			".mp3", ".aac", ".m4a", ".flac", ".ogg", ".wav",
		}))
		fd.Show()
	}

	addFolderBtn.OnTapped = func() {
		mu.Lock()
		mp := currentMountpoint
		mu.Unlock()
		if mp == "" {
			dialog.ShowError(fmt.Errorf("no iPod connected"), w)
			return
		}
		fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			var paths []string
			filepath.WalkDir(uri.Path(), func(p string, d os.DirEntry, err error) error {
				if err == nil && !d.IsDir() && isAudio(p) {
					paths = append(paths, p)
				}
				return nil
			})
			if len(paths) == 0 {
				dialog.ShowInformation("No audio files", "No supported audio files found.", w)
				return
			}
			dialog.ShowConfirm("Add folder",
				fmt.Sprintf("Add %d file(s) to iPod?", len(paths)),
				func(ok bool) {
					if !ok {
						return
					}
					setBusy(true)
					go runBatch(mp, paths, "Adding", func(p string) error {
						return doAddTrack(mp, p)
					})
				}, w)
		}, w)
		fd.Show()
	}

	deleteBtn.OnTapped = func() {
		mu.Lock()
		mp := currentMountpoint
		mu.Unlock()
		if mp == "" {
			dialog.ShowError(fmt.Errorf("no iPod connected"), w)
			return
		}
		type item struct {
			id    int
			title string
		}
		var items []item
		for i, checked := range selected {
			if checked && i < len(filteredIndices) {
				t := fullTracks[filteredIndices[i]]
				items = append(items, item{t.ID, t.Title})
			}
		}
		if len(items) == 0 {
			dialog.ShowInformation("Nothing selected", "Check at least one song to delete.", w)
			return
		}
		dialog.ShowConfirm("Delete songs",
			fmt.Sprintf("Permanently delete %d song(s)?", len(items)),
			func(ok bool) {
				if !ok {
					return
				}
				// convert to paths for runBatch (reuse title as label)
				titles := make([]string, len(items))
				ids    := make([]int, len(items))
				for i, it := range items {
					titles[i] = it.title
					ids[i]    = it.id
				}
				setBusy(true)
				go func() {
					total := len(ids)
					var errs []string
					for i, id := range ids {
						idx := i + 1
						fyne.Do(func() {
							progressLabel.SetText(fmt.Sprintf("Deleting %d / %d", idx, total))
						})
						if err := doDeleteTrack(mp, id); err != nil {
							errs = append(errs, fmt.Sprintf("%s: %v", titles[i], err))
						}
					}
					go refreshStorage(mp)
					refreshTracks(mp)
					fyne.Do(func() {
						setBusy(false)
						if len(errs) > 0 {
							dialog.ShowError(fmt.Errorf("%s", strings.Join(errs, "\n")), w)
						}
					})
				}()
			}, w)
	}

	selectAllBtn := widget.NewButton("Select All", func() {
		for i := range filteredIndices {
			selected[i] = true
		}
		songList.Refresh()
	})

	selectNoneBtn := widget.NewButton("Select None", func() {
		for k := range selected {
			delete(selected, k)
		}
		songList.Refresh()
	})

	infoForm := widget.NewForm(
		widget.NewFormItem("Model",    modelVal),
		widget.NewFormItem("GUID",     guidVal),
		widget.NewFormItem("Capacity", capacityVal),
		widget.NewFormItem("Used",     usedVal),
		widget.NewFormItem("Free",     freeVal),
		widget.NewFormItem("Songs",    songsVal),
	)

	leftPanel := container.NewBorder(
		container.NewVBox(statusLabel, toggleBtn, widget.NewSeparator(), infoForm),
		container.NewPadded(container.NewVBox(addBtn, addFolderBtn)),
		nil, nil,
	)

	rightPanel := container.NewBorder(
		searchEntry,
		container.NewPadded(container.NewHBox(
			deleteBtn, selectAllBtn, selectNoneBtn,
			widget.NewLabel(""),
			progressLabel,
		)),
		nil, nil,
		songList,
	)

	hsplit := container.NewHSplit(leftPanel, rightPanel)
	hsplit.SetOffset(0.30)
	w.SetContent(hsplit)

	go func() {
		mp := detectMountpoint()
		if mp == "" {
			return
		}
		info := fetchDeviceInfo(mp)
		if !info.Connected {
			return
		}
		fyne.Do(func() {
			mu.Lock()
			currentMountpoint = mp
			mu.Unlock()
			updateDeviceUI(info)
			toggleBtn.SetText("Disconnect")
			setBusy(true)
		})
		refreshTracks(mp)
		fyne.Do(func() { setBusy(false) })
	}()

	w.ShowAndRun()
}
