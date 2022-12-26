/* placeholder file for JavaScript */
const confirm_delete = (id) => {
    if(window.confirm(`Task ${id} を削除しますか？`)) {
        location.href = `/task/delete/${id}`;
    }
}
 
const confirm_update = (id) => {
    if(window.confirm(`Task ${id} を編集しますか？`)) {
        location.href = `/task/edit/${id}`;
    }
}