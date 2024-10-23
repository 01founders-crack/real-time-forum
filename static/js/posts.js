export async function fetchAndDisplayPosts() {
    try {
        const response = await fetch('http://localhost:8080/posts');
        const posts = await response.json();
        const container = document.getElementById('feed-container');

        posts.forEach(post => {
            const card = document.createElement('div');
            card.className = 'post-card';
            if (`${post.category}` != '') {
                card.classList.add(`${post.category}`)
            }
            

                const details = document.createElement('div');
                details.className = 'post-details';

                    const user = document.createElement('p');
                    user.textContent = `${post.username}`;
                    details.appendChild(user);

                const postMain = document.createElement('div');
                postMain.className = 'post-main';

                    const postTitle = document.createElement('h2');
                    postTitle.textContent = `${post.title}`;
                    postTitle.className = 'postTitle'
                    postMain.appendChild(postTitle);

                    const postBody = document.createElement('p')
                    postBody.textContent = `${post.content}`;
                    postBody.className = 'postBody'
                    postMain.appendChild(postBody);

                const commentBlock = document.createElement('div');
                    commentBlock.className = 'commentBlock'
                    const commentShowHide = document.createElement('p')
                        commentShowHide.className = 'showHide'
                        commentShowHide.textContent = 'Comments'
                    // const addComment = document.createElement('p')
                    //     addComment.className = 'addComment'
                    //     addComment.textContent = 'Add Comment'
                    const addCommentForm = document.createElement('div')
                    addCommentForm.className = 'commentsForm'
                    const createComment = document.createElement('form')
                        createComment.action = '/addComment'
                        createComment.method = 'POST'
                        createComment.className = 'createComment'
                        const postToComment = document.createElement('input')
                        postToComment.name = 'postId'
                        postToComment.value = `${post.stringId}`
                        postToComment.style.display = 'none'
                        const commentMessage = document.createElement('textarea')
                            commentMessage.rows = ''
                            commentMessage.cols = '30'
                            commentMessage.name = 'commentMessage'
                            commentMessage.className = 'commentMessage'
                            commentMessage.placeholder = 'Write comment here.'
                        const submitComment = document.createElement('button')
                            submitComment.type = 'submit'
                            submitComment.textContent = 'Add Comment'
                        createComment.appendChild(postToComment)
                        createComment.appendChild(commentMessage)
                        createComment.appendChild(submitComment)
                    addCommentForm.appendChild(createComment)
                commentBlock.appendChild(commentShowHide)
                commentBlock.appendChild(addCommentForm)
                    const comments = document.createElement('div');
                    comments.className = 'comments'
                
                    if (post.comments != null) {
                        post.comments.forEach(singleComment => {
                            const comment = document.createElement('div');
                            comment.className = 'comment';

                                const commenter = document.createElement('p');
                                commenter.textContent = `${singleComment.username}`;
                                commenter.className = 'commenter'
                                comment.appendChild(commenter);

                                const commentBody = document.createElement('p');
                                commentBody.textContent = `${singleComment.content}`;
                                comment.appendChild(commentBody);

                            comments.appendChild(comment);
                            
                        });
                    }
                    commentBlock.appendChild(comments)

                card.appendChild(details);
                card.appendChild(postMain);
                card.appendChild(commentBlock);
                
                
                container.appendChild(card);
                
                
            
        });
    } catch(error) {
        console.error('Error fetching posts:', error);
    }

    
}
