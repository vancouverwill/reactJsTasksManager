<?php
/**
 * This file provided by Facebook is for non-commercial testing and evaluation
 * purposes only. Facebook reserves all rights not expressly granted.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * FACEBOOK BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
$scriptInvokedFromCli =
    isset($_SERVER['argv'][0]) && $_SERVER['argv'][0] === 'server.php';

if($scriptInvokedFromCli) {
    $port = getenv('PORT');
    if (empty($port)) {
        $port = "3000";
    }

    echo 'starting server on port '. $port . PHP_EOL;
    exec('php -S localhost:'. $port . ' -t public server.php');
} else {
    return routeRequest();
}

function routeRequest()
{
    if (file_exists('tasks.json')) {
        $tasks = file_get_contents('tasks.json');
    }
    else {
        $tasks = "[]";
    }

    $uri = $_SERVER['REQUEST_URI'];
    if ($uri == '/') {
        echo file_get_contents('./public/index.html');
    } elseif (preg_match('/\/api\/tasks(\?.*)?/', $uri)) {
        if($_SERVER['REQUEST_METHOD'] === 'POST') {
            $tasksDecoded = json_decode($tasks, true);
            $tasksDecoded[] = ['id' => count($tasksDecoded) + 1, 'content' => $_POST['content'], 'other' => $_POST['other']];

                $tasks = json_encode($tasksDecoded, JSON_PRETTY_PRINT);
                file_put_contents('tasks.json', $tasks);
            }
            elseif($_SERVER['REQUEST_METHOD'] === 'DELETE') {
                $tasksDecoded = json_decode($tasks, true);
                $content = file_get_contents('php://input');

                $array = array();

                parse_str($content, $array);

                if (isset($array["id"])) {

                    foreach($tasksDecoded AS $index => $task) {
                        if ($task["id"] == $array["id"]) {
                            unset($tasksDecoded[$index]);
                        }
                    }
                    $tasks = json_encode(array_values($tasksDecoded), JSON_PRETTY_PRINT);

                    file_put_contents('tasks.json', $tasks);
                }
            }


        header('Content-Type: application/json');
        header('Cache-Control: no-cache');
        echo $tasks;
    } else {
        return false;
    }
}
