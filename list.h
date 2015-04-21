#ifndef LIST_H
#define LIST_H

#include <stdlib.h>
#include <assert.h>
#include <string.h>

typedef struct L {
    struct L* prev;
    struct L* next;
    void* item;
} ListNode;

ListNode* node_alloc();
ListNode* node_init(ListNode* prev, ListNode* next, void* item);
void node_free(ListNode* node);
void node_remove(ListNode* node);
ListNode* node_insert_after(ListNode* node, void* item);
void* node_get_item(ListNode* node, int index);
ListNode* node_find_item(ListNode* node, void* item);

typedef struct {
    ListNode* head;
    ListNode* tail;
} List;

List* list_init();
ListNode* list_append(List* list, void* item);
int  list_remove(List* list, void* item);
size_t list_len(List* list);
void list_free();

#endif
