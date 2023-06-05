
function createPostMethod(className, path){
    let links = document.getElementsByClassName(className);

    for (let i = 0; i < links.length; i++) {
        links[i].addEventListener('click', function(e) {
            e.preventDefault();
            
            let id = this.getAttribute('data-id');
            let roundId = this.getAttribute('data-round');
            
            // Создаем форму для отправки POST-запроса
            let form = document.createElement('form');
            form.method = 'POST';
            form.action = path

            // Создаем поле для передачи ID
            let idField = document.createElement('input');
            idField.type = 'hidden';
            idField.name = 'id';
            idField.value = id;

            let roundField = document.createElement('input');
            roundField.type = 'hidden';
            roundField.name = 'round_number';
            roundField.value = roundId;
            
            // Добавляем поле в форму и добавляем форму на страницу
            form.appendChild(idField);
            form.appendChild(roundField);
            document.body.appendChild(form);
            
            // Отправляем форму
            form.submit();
        });
    }
}