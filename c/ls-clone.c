#include "stdio.h"
#include "dirent.h"
#include "stdlib.h"
#include "string.h"
#include <sys/stat.h>
#include "stdbool.h"

#define MAX_PATH_LEN 1000

struct {
    unsigned int a: 1;
    unsigned int l: 1;
} flags;

bool is_file(const char* path);
bool is_dir(const char* path);
int print_dir(const char* path);
struct stat get_file_info(const char* path);
void parse_args(char *path, int argc, char *argv[]);

int main(int argc, char *argv[]) {
    // default to current directory
    char path[MAX_PATH_LEN] = ".";

    parse_args(path, argc, argv);

    if(is_file(path)){
        struct stat info = get_file_info(path);
        printf("%lld\n", info.st_size);
    } else {
        print_dir(path);
    }
}

void parse_args(char *path, int argc, char *argv[]){
    for(int i = 1; i < argc; i++){
        if(argv[i][0] == '-'){
            for(int j = 1; j < sizeof argv[i]; j++){
                switch (argv[i][j]) {
                    case 'a':
                        flags.a = 1;
                    case 'l':
                        flags.l = 1;
                }
            }
        } else {
            strcpy(path, argv[i]);
        }
    }
}


bool is_file(const char* path){
    struct stat buf;
    stat(path, &buf);
    return S_ISREG(buf.st_mode);
}

bool is_dir(const char* path){
    struct stat buf;
    stat(path, &buf);
    return S_ISDIR(buf.st_mode);
}

int print_dir(const char* path) {
    DIR *dir;
    struct dirent *ent;
    if ((dir = opendir(path)) != NULL) {
        // print all the files and directories within directory
        while ((ent = readdir (dir)) != NULL) {
            if(!flags.a && (ent->d_name)[0] == '.')
                ;
            else
                // TODO: output whatever info has been requested
                printf ("%s\n", ent->d_name);
        }
        closedir (dir);
        return EXIT_SUCCESS;
    } else {
        // could not open directory
        perror ("");
        return 1;
    }
}

struct stat get_file_info(const char* path){
    struct stat buf;
    stat(path, &buf);
    return buf;
}