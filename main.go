package main

// #cgo LDFLAGS: -lrsync
// #include <stdio.h>
// #include <stdlib.h>
// #include <librsync.h>
import "C"
import (
	"log"
	"unsafe"
)

func Signature(oldFile string, sigFile string) {
	cOldFile := C.CString(oldFile)
	cSigFile := C.CString(sigFile)
	cFopenRb := C.CString("rb")
	cFopenWb := C.CString("wb")
	defer C.free(unsafe.Pointer(cOldFile))
	defer C.free(unsafe.Pointer(cSigFile))
	defer C.free(unsafe.Pointer(cFopenRb))
	defer C.free(unsafe.Pointer(cFopenWb))

	basis := C.fopen(cOldFile, cFopenRb)
	signature := C.fopen(cSigFile, cFopenWb)
	defer C.fclose(basis)
	defer C.fclose(signature)

	var stats C.struct_rs_stats
	var result C.rs_result = C.rs_sig_file(basis, signature, C.RS_DEFAULT_BLOCK_LEN, C.RS_DEFAULT_STRONG_LEN, &stats)

	log.Println(result)
}

func Delta(newFilePath string, sigFilePath string, deltaFilePath string) {
	cNewFilePath := C.CString(newFilePath)
	cSigFilePath := C.CString(sigFilePath)
	cDeltaFilePath := C.CString(deltaFilePath)
	cFopenRb := C.CString("rb")
	cFopenWb := C.CString("wb")
	defer C.free(unsafe.Pointer(cNewFilePath))
	defer C.free(unsafe.Pointer(cSigFilePath))
	defer C.free(unsafe.Pointer(cDeltaFilePath))
	defer C.free(unsafe.Pointer(cFopenRb))
	defer C.free(unsafe.Pointer(cFopenWb))

	newFile := C.fopen(cNewFilePath, cFopenRb)
	sigFile := C.fopen(cSigFilePath, cFopenRb)
	deltaFile := C.fopen(cDeltaFilePath, cFopenWb)
	defer C.fclose(newFile)
	defer C.fclose(sigFile)
	defer C.fclose(deltaFile)

	var result C.rs_result
	var stats C.rs_stats_t
	var signature *C.rs_signature_t

	result = C.rs_loadsig_file(sigFile, &signature, &stats)
	if result != C.RS_DONE {
		log.Println("Fail when load Signature File")
		return
	}

	result = C.rs_build_hash_table(signature)
	if result != C.RS_DONE {
		log.Println("Fail when try to index Signature")
		return
	}

	result = C.rs_delta_file(signature, newFile, deltaFile, &stats)
	defer C.rs_free_sumset(signature)
}

func Patch(oldFilePath string, deltaFilePath string, patchedFilePath string) {
	cOldFilePath := C.CString(oldFilePath)
	cDeltaFilePath := C.CString(deltaFilePath)
	cPatchedFilePath := C.CString(patchedFilePath)
	cFopenRb := C.CString("rb")
	cFopenWb := C.CString("wb")
	defer C.free(unsafe.Pointer(cOldFilePath))
	defer C.free(unsafe.Pointer(cDeltaFilePath))
	defer C.free(unsafe.Pointer(cPatchedFilePath))
	defer C.free(unsafe.Pointer(cFopenRb))
	defer C.free(unsafe.Pointer(cFopenWb))

	oldFile := C.fopen(cOldFilePath, cFopenRb)
	deltaFile := C.fopen(cDeltaFilePath, cFopenRb)
	patchedFile := C.fopen(cPatchedFilePath, cFopenWb)
	defer C.fclose(oldFile)
	defer C.fclose(deltaFile)
	defer C.fclose(patchedFile)

	var result C.rs_result
	var stats C.rs_stats_t

	result = C.rs_patch_file(oldFile, deltaFile, patchedFile, &stats)
	log.Println(result)
}

func main() {
	log.Println("synko loaded..")
}
