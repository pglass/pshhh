#include "list.h"
#include <stdio.h>

ListNode* node_alloc() {
    ListNode* node = (ListNode*) malloc(sizeof(ListNode));
    memset((void*) node, 0, sizeof(ListNode));
    return node;
}

ListNode* node_init(ListNode* prev, ListNode* next, void* item) {
    ListNode* node = node_alloc();
    node->prev = prev;
    node->next = next;
    node->item = item;
    return node;
}

void node_free(ListNode* node) {
    free(node);
}

ListNode* node_insert_after(ListNode* node, void* item) {
    ListNode* new = node_init(node, node->next, item);
    node->next = new;
    return new;
}

void node_remove(ListNode* node) {
    assert(node != NULL);
    if (node->prev != NULL & node->next != NULL) {
        node->prev->next = node->next;
        node->next->prev = node->prev;
    } else if (node->prev == NULL) {
        node->next->prev = NULL;
    } else if (node->next == NULL) {
        node->prev->next = NULL;
    }
}

void* node_get_item(ListNode* node, int index) {
    int i = 0;
    while (node != NULL) {
        if (i == index) {
            return node->item;
        }
        i++;
        node = node->next;
    }
    return NULL;
}

ListNode* node_find_item(ListNode* node, void* item) {
    while (node != NULL) {
        if (node->item == item) {
            return node;
        }
        node = node->next;
    }
    return NULL;
}

List* list_init() {
    List* list = (List*) malloc(sizeof(List));
    list->head = node_alloc();
    list->tail = node_alloc();
    list->tail->prev = list->head;
    list->head->next = list->tail;
    return list;
}

ListNode* list_append(List* list, void* item) {
    ListNode* new = node_init(list->tail->prev, list->tail->prev->next, item);
    list->tail->prev->next = new;
    list->tail->prev = new;
    assert(list->tail->next == NULL);
    assert(list->head->prev == NULL);
    assert(list->tail->prev == new);
    return new;
}

int list_remove(List* list, void* item) {
    ListNode* node = node_find_item(list->head, item);
    if (node != NULL) {
        node_remove(node);
        node_free(node);
        return 1;
    }
    return 0;
}

/* Does not call free on each item */
void list_free(List* list) {
    ListNode* node = list->head;
    while (node != NULL) {
        ListNode* next = node->next;
        printf("free %p\n", node);
        node_free(node);
        node = next;
    }
    free(list);
}
