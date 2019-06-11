document.addEventListener("DOMContentLoaded", function () {
    FormHandler.handleSubmitButtons();

    let tabHandler = new TabHandler();
    tabHandler.addEventListener('.tabs.tabs-entry-edit', EntryEditorHandler.switchHandler);

    let navHandler = new NavHandler();

    if (! document.querySelector("body[data-disable-keyboard-shortcuts=true]")) {
        let keyboardHandler = new KeyboardHandler();
        keyboardHandler.on("g u", () => navHandler.goToPage("unread"));
        keyboardHandler.on("g b", () => navHandler.goToPage("starred"));
        keyboardHandler.on("g h", () => navHandler.goToPage("history"));
        keyboardHandler.on("g f", () => navHandler.goToFeedOrFeeds());
        keyboardHandler.on("g c", () => navHandler.goToPage("categories"));
        keyboardHandler.on("g s", () => navHandler.goToPage("settings"));
        keyboardHandler.on("ArrowLeft", () => navHandler.goToPrevious());
        keyboardHandler.on("ArrowRight", () => navHandler.goToNext());
        keyboardHandler.on("k", () => navHandler.goToPrevious());        
        keyboardHandler.on("p", () => navHandler.goToPrevious());
        keyboardHandler.on("j", () => navHandler.goToNext());
        keyboardHandler.on("n", () => navHandler.goToNext());
        keyboardHandler.on("h", () => navHandler.goToPage("previous"));
        keyboardHandler.on("l", () => navHandler.goToPage("next"));
        keyboardHandler.on("o", () => navHandler.openSelectedItem());
        keyboardHandler.on("v", () => navHandler.openOriginalLink());
        keyboardHandler.on("m", () => navHandler.toggleEntryStatus());
        keyboardHandler.on("A", () => {
            let element = document.querySelector("a[data-on-click=markPageAsRead]");
            navHandler.markPageAsRead(element.dataset.showOnlyUnread || false);
        });
        keyboardHandler.on("s", () => navHandler.saveEntry());
        keyboardHandler.on("d", () => navHandler.fetchOriginalContent());
        keyboardHandler.on("f", () => navHandler.toggleBookmark());
        keyboardHandler.on("?", () => navHandler.showKeyboardShortcuts());
        keyboardHandler.on("#", () => navHandler.unsubscribeFromFeed());
        keyboardHandler.on("/", (e) => navHandler.setFocusToSearchInput(e));
        keyboardHandler.on("Escape", () => ModalHandler.close());
        keyboardHandler.listen();
    }

    let touchHandler = new TouchHandler(navHandler);
    touchHandler.listen();

    let mouseHandler = new MouseHandler();
    mouseHandler.onClick("a[data-save-entry]", (event) => {
        EntryHandler.saveEntry(event.target);
    });

    mouseHandler.onClick("a[data-toggle-bookmark]", (event) => {
        EntryHandler.toggleBookmark(event.target);
    });

    mouseHandler.onClick("a[data-toggle-cache]", (event) => {
        EntryHandler.toggleCache(event.target);
    });

    mouseHandler.onClick("a[data-history-go-back]", (event) => {
        history.go(-1);
    });

    mouseHandler.onClick("a[data-toggle-status]", (event) => {
        let currentItem = DomHelper.findParent(event.target, "entry");
        if (!currentItem) {
            currentItem = DomHelper.findParent(event.target, "item");
        }

        if (currentItem) {
            EntryHandler.toggleEntryStatus(currentItem);
        }
    });

    mouseHandler.onClick("a[data-set-read]", (event) => {
        let currentItem = DomHelper.findParent(event.target, "entry");
        if (!currentItem) {
            currentItem = DomHelper.findParent(event.target, "item");
        }
        if (currentItem) {
            EntryHandler.setEntryStatusRead(currentItem)
        }

    }, true);

    mouseHandler.onClick("a[data-fetch-content-entry]", (event) => {
        EntryHandler.fetchOriginalContent(event.target);
    });

    mouseHandler.onClick("a[data-on-click=showActionMenu]", (event) => {
        let currentItem = DomHelper.findParent(event.target, "entry");
        if (!currentItem) {
            currentItem = DomHelper.findParent(event.target, "item");
        }
        if (currentItem) {
            new ActionMenuHandler(currentItem).show();
        }
    })

    mouseHandler.onClick("a[data-on-click=markPageAsRead]", (event) => {
        navHandler.markPageAsRead(event.target.dataset.showOnlyUnread || false);
    });

    mouseHandler.onClick("a[data-confirm]", (event) => {
        (new ConfirmHandler()).handle(event);
    });

    mouseHandler.onClick("a[data-action=search]", (event) => {
        navHandler.setFocusToSearchInput(event);
    });

    mouseHandler.onClick("button[data-action=submit-entry]", (event) => {
        EntryEditorHandler.submitHandler(event);
    });

    mouseHandler.onClick("a[data-link-state=flip]", (event) => {
        LinkStateHandler.flip(event.target);
    }, true);

    let menuHandler = new MenuHandler();
    mouseHandler.onClick(".logo", (event) => menuHandler.logoClickHandler(event));

    if ("serviceWorker" in navigator) {
        let scriptElement = document.getElementById("service-worker-script");
        if (scriptElement) {
            navigator.serviceWorker.register(scriptElement.src);
        }
    }
});

window.onload = function () {
    // masonry has to wait for all resources loaded to get the right layout
    let msnryElement = document.querySelector('.masonry');
    if (msnryElement) {
        var msnry = new Masonry(msnryElement, {
            itemSelector: '.item',
            columnWidth: '.item-sizer',
            gutter: 10
        })
        let callback = (instance, image) => {
            if (image && image.img && !image.img.dataset.src && !image.isLoaded) {
                let thumbnail = DomHelper.findParent(image.img, "thumbnail");
                if (thumbnail) thumbnail.parentNode.removeChild(thumbnail);
            }
            msnry.layout();
        }
        imagesLoaded('.masonry .item').on('progress', callback);
        LazyloadHandler.add(".item", 'progress', callback);
    }
};